package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type Service struct {
}

var Validate = validator.New()

func NewService() *Service {
	return &Service{}
}

type DataTopic struct {
	Topic    string `json:"topic"`
	Datatype string `json:"datatype"`
}

type DeviceInfo struct {
	Topic       []string `json:"topic"`
	Broker      string   `json:"broker"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Device_type string   `json:"device_type"`
}

type Error struct {
	ApiVersion string `json:"apiVersion"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func (s *Service) AddDevice(w http.ResponseWriter, r *http.Request) {

	var deviceInfo DeviceInfo

	if err := ParseJSON(r, &deviceInfo); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := Validate.Struct(deviceInfo); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := validateFields(deviceInfo); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	file, errFile := os.Open("../docker-compose.yml")

	if errFile != nil {
		WriteError(w, http.StatusBadRequest, errFile)
		return
	}
	defer file.Close()

	var search_string = fmt.Sprintf("device-mqtt-broker-%s:", deviceInfo.Broker) // Replace with actual broker name
	fmt.Println(search_string)

	var port int

	scanner := bufio.NewScanner(file)

	found := false
	re := regexp.MustCompile(`59982:(\d+)`)
	var maxPort float64 = math.Inf(-1)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			// matches[1] contains the first captured group (the first number after 127.0.0.1:)
			fmt.Println("Found number:", matches[1])
			portStr := matches[1]              // Port number as string
			port, err := strconv.Atoi(portStr) // Convert the port string to an integer
			if err != nil {
				fmt.Println("Error converting port:", err)
				continue
			}
			maxPort = math.Max(maxPort, float64(port)+1)
		}
		// Check if the line contains the search string
		if strings.Contains(line, search_string) {
			found = true
			break // Exit loop once found
		}
	}
	file.Close()
	var portString string
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	if maxPort == math.Inf(-1) {
		fmt.Println("No matching ports found.")
	} else {
		port = int(maxPort) + 1
		fmt.Println("Maximum port found:", int(maxPort)) // Convert maxPort back to int for printing
		portString = strconv.Itoa(port)

	}

	fmt.Println(port, found)
	if !found {
		var formattedData = fmt.Sprintf(`

  device-mqtt-broker-%s:
    command: /device-mqtt -cp=consul.http://edgex-core-consul:8500 --registry --confdir=/res
    container_name: edgex-device-mqtt-broker-%s
    depends_on:
    - consul
    - data
    - metadata
    - security-bootstrapper
    entrypoint:
    - /edgex-init/ready_to_run_wait_install.sh
    environment:
      MQTTBROKERINFO_HOST: %s
      MQTTBROKERINFO_PORT: 1883
      MQTTBROKERINFO_QOS: 0
      MQTTBROKERINFO_AUTHMODE: usernamepassword
      MQTTBROKERINFO_USERNAME: %s
      MQTTBROKERINFO_PASSWORD: %s
      # change the client ID to resolve conflict if there are more than 1 machine running the service at the same time
      MQTTBROKERINFO_CLIENTID: mqtt-client-deploy-AIOT
      MQTTBROKERINFO_INCOMINGTOPIC: STP/#
      MQTTBROKERINFO_RESPONSETOPIC: STP/#
      MQTTBROKERINFO_USETOPICLEVELS: "true"
      API_GATEWAY_HOST: edgex-kong
      API_GATEWAY_STATUS_PORT: '8100'
      CLIENTS_CORE_COMMAND_HOST: edgex-core-command
      CLIENTS_CORE_DATA_HOST: edgex-core-data
      CLIENTS_CORE_METADATA_HOST: edgex-core-metadata
      CLIENTS_SUPPORT_NOTIFICATIONS_HOST: edgex-support-notifications
      CLIENTS_SUPPORT_SCHEDULER_HOST: edgex-support-scheduler
      DATABASES_PRIMARY_HOST: edgex-redis
      EDGEX_SECURITY_SECRET_STORE: "true"
      MESSAGEQUEUE_HOST: edgex-redis
      PROXY_SETUP_HOST: edgex-security-proxy-setup
      REGISTRY_HOST: edgex-core-consul
      SECRETSTORE_HOST: edgex-vault
      SECRETSTORE_PORT: '8200'
      SERVICE_HOST: edgex-device-mqtt-broker-%s
      SPIFFE_ENDPOINTSOCKET: /tmp/edgex/secrets/spiffe/public/api.sock
      SPIFFE_TRUSTBUNDLE_PATH: /tmp/edgex/secrets/spiffe/trust/bundle
      SPIFFE_TRUSTDOMAIN: edgexfoundry.org
      STAGEGATE_BOOTSTRAPPER_HOST: edgex-security-bootstrapper
      STAGEGATE_BOOTSTRAPPER_STARTPORT: '54321'
      STAGEGATE_DATABASE_HOST: edgex-redis
      STAGEGATE_DATABASE_PORT: '6379'
      STAGEGATE_DATABASE_READYPORT: '6379'
      STAGEGATE_KONGDB_HOST: edgex-kong-db
      STAGEGATE_KONGDB_PORT: '5432'
      STAGEGATE_KONGDB_READYPORT: '54325'
      STAGEGATE_READY_TORUNPORT: '54329'
      STAGEGATE_REGISTRY_HOST: edgex-core-consul
      STAGEGATE_REGISTRY_PORT: '8500'
      STAGEGATE_REGISTRY_READYPORT: '54324'
      STAGEGATE_SECRETSTORESETUP_HOST: edgex-security-secretstore-setup
      STAGEGATE_SECRETSTORESETUP_TOKENS_READYPORT: '54322'
      STAGEGATE_WAITFOR_TIMEOUT: 60s
      # WRITETABLE_INSECURESECRETS_MQTT_PATH: "credentials"
      # WRITETABLE_INSECURESECRETS_MQTT_SECRETS_USERNAME: "username"
      # WRITETABLE_INSECURESECRETS_MQTT_SECRETS_PASSWORD: "password"
    hostname: edgex-device-mqtt-broker-%s
    image: docker.io/edgexfoundry/device-mqtt:0.0.0-dev
    networks:
      edgex-network: {}
    ports:
    - 59982:%s/tcp
    read_only: true
    restart: always
    security_opt:
    - no-new-privileges:true
    user: 2002:2001
    volumes:
    - edgex-init:/edgex-init:ro,z
    - /tmp/edgex/secrets/device-mqtt-broker-%s:/tmp/edgex/secrets/device-mqtt-broker-%s:ro,z`,
			deviceInfo.Broker, deviceInfo.Broker, deviceInfo.Broker, deviceInfo.Username, deviceInfo.Password, deviceInfo.Broker, deviceInfo.Broker, portString)
		// , deviceInfo.Broker, deviceInfo.Broker, deviceInfo.Broker, deviceInfo.Broker, deviceInfo.Broker, port)

		file, errFile := os.OpenFile("../docker-compose.yml", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if errFile != nil {
			fmt.Println(errFile)
			return
		}
		defer file.Close()

		if _, err := file.Write([]byte(formattedData)); err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}
		if err := scanAndUpdate("../docker-compose.yml", fmt.Sprintf("device-mqtt-broker-%s", deviceInfo.Broker)); err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	// fmt.Println(deviceInfo.Broker)

	// url_edgex := fmt.Sprintf("http://localhost:59881/deviceservice/name/device-mqtt", deviceInfo.Broker)
	url_edgex := fmt.Sprintf("http://localhost:59881/api/v2/deviceservice/name/device-mqtt")

	c := http.Client{Timeout: time.Duration(1) * time.Second}
	// fmt.Println("Here")
	if !getRequest(c, url_edgex) {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("Error connect to device service"))
	}

	var device_name []string
	var profile_name []string
	var todo []string
	success := false
	fmt.Println(len(deviceInfo.Topic))
	for i := 0; i < len(deviceInfo.Topic); i++ {
		parts := strings.Split(deviceInfo.Topic[i], "/")
		if len(parts) >= 3 {
			device_name = append(device_name, parts[len(parts)-2])
			if parts[len(parts)-1] != "data" && parts[len(parts)-1] != "status" {
				WriteError(w, http.StatusBadRequest, fmt.Errorf("Profile name is not valid"))
				return
			}
			profile_name = append(profile_name, parts[len(parts)-1])

		} else {
			WriteError(w, http.StatusBadRequest, fmt.Errorf("Topic is not valid"))
			return
		}
	}

	for i := 0; i < len(deviceInfo.Topic); i++ {

		var url_device_profile = fmt.Sprintf("http://localhost:59881/api/v2/device/name/%s", device_name[i])
		var url_device_resource = fmt.Sprintf("http://localhost:59881/api/v2/deviceresource/profile/%s-Profile/resource/%s", device_name[i], profile_name[i])
		// Send the GET request

		if getRequest(c, url_device_profile) {
			if !getRequest(c, url_device_resource) {
				todo = append(todo, "update")
			} else {
				WriteError(w, http.StatusBadRequest, fmt.Errorf("Device %s already exist", device_name[i]))
				return
			}
		} else {
			// else {
			todo = append(todo, "post")
			// }
		}
	}

	url_profile := "http://localhost:59881/api/v2/deviceprofile"
	url_device := "http://localhost:59881/api/v2/device"
	url_resource := "http://localhost:59881/api/v2/deviceprofile/resource"
	for i := 0; i < len(deviceInfo.Topic); i++ {

		if todo[i] == "post" {
			data := []map[string]interface{}{
				{
					"apiVersion": "v2",
					"profile": map[string]interface{}{
						"name":         fmt.Sprintf("%s-Profile", device_name[i]),
						"manufacturer": "SMIC",
						"model":        "1",
						"labels": []string{
							"MQTT",
							"data",
						},
						"description": "device profile of MQTT devices",
						"deviceResources": []map[string]interface{}{
							{
								"name":        profile_name[i],
								"isHidden":    false,
								"description": "data JSON message",
								"properties": map[string]interface{}{
									"valueType": "Object",
									"readWrite": "RW",
									"mediaType": "application/json",
								},
							},
						},
					},
				},
			}

			jsonDataProfile, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				WriteError(w, http.StatusBadRequest, err)
				return
			}
			if err := postRequest(c, url_profile, jsonDataProfile); err != nil {
				WriteError(w, http.StatusBadRequest, err)
				return
			} else {

				device := []map[string]interface{}{{
					"apiVersion": "v2",
					"device": map[string]interface{}{
						"name":           device_name[i],
						"description":    "Test mqtt",
						"adminState":     "UNLOCKED",
						"operatingState": "UP",
						"labels": []string{
							"home", "mqtt",
						},
						"serviceName": "device-mqtt",
						"profileName": fmt.Sprintf("%s-Profile", device_name[i]),
						"protocols": map[string]interface{}{
							"mqtt": map[string]string{
								"CommandTopic": deviceInfo.Topic[i],
							},
						},
					},
				}}

				jsonDataProfile, err := json.MarshalIndent(device, "", "  ")
				if err != nil {
					fmt.Println("Error marshaling JSON:", err)
					WriteError(w, http.StatusBadRequest, err)
					return
				}
				if err := postRequest(c, url_device, jsonDataProfile); err != nil {
					WriteError(w, http.StatusBadRequest, err)
					return
				} else {
					success = true

				}

			}
		} else if todo[i] == "update" {
			fmt.Println("Ready to update")
			data := []map[string]interface{}{
				{
					"apiVersion":  "v2",
					"profileName": fmt.Sprintf("%s-Profile", device_name[i]),
					"resource": map[string]interface{}{
						"description": "data JSON message",
						"isHidden":    false,
						"name":        profile_name[i],
						"properties": map[string]interface{}{
							"mediaType": "application/json",
							"readWrite": "RW",
							"valueType": "Object",
						},
					},
				},
			}
			jsonDataProfile, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				WriteError(w, http.StatusBadRequest, err)
				return
			}
			fmt.Println(string(jsonDataProfile))
			if err := postRequest(c, url_resource, jsonDataProfile); err != nil {
				WriteError(w, http.StatusBadRequest, err)
				return
			} else {
				success = true
			}
		} else {
			WriteError(w, http.StatusBadRequest, fmt.Errorf("Unexpected error"))
			return
		}
	}

	if success {
		WriteJSON(w, http.StatusOK, fmt.Sprintf("Add and Update Device Success"))
		return
	}

}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	// fmt.Println(r.Body)

	return json.NewDecoder(r.Body).Decode(v)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func validateFields(deviceInfo DeviceInfo) error {
	if len(deviceInfo.Topic) == 0 {
		return fmt.Errorf("topic field is required")
	}
	if deviceInfo.Broker == "" {
		return fmt.Errorf("broker field is required")
	}
	if deviceInfo.Username == "" {
		return fmt.Errorf("username field is required")
	}
	if deviceInfo.Password == "" {
		return fmt.Errorf("password field is required")
	}
	return nil
}
