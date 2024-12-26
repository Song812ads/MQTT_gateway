package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	// file, errFile := os.OpenFile("../docker-compose.yml", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	// if errFile != nil {
	// 	WriteError(w, http.StatusBadRequest, errFile)
	// 	return
	// }
	// defer file.Close()

	// var search_string = fmt.Sprintf("device-mqtt-%s", deviceInfo.Broker) // Replace with actual broker name

	// fmt.Println(search_string)

	// scanner := bufio.NewScanner(file)

	// found := false
	// // re := regexp.MustCompile(`127\.0\.0\.1:(\d+):59982/tcp`)

	// for scanner.Scan() {
	// 	line := scanner.Text()

	// 	// Check if the line contains the search string
	// 	if strings.Contains(line, search_string) {
	// 		found = true
	// 		fmt.Println("Exist")
	// 		break // Exit loop once found
	// 	}
	// }

	// if !found {

	// 	// for scanner.Scan() {
	// 	// 	line := scanner.Text()
	// 	// 	matches := re.FindStringSubmatch(line)
	// 	// 	// if len(matches) > 1 {
	// 	// 	// matches[1] contains the first captured group (the first number after 127.0.0.1:)
	// 	// 	fmt.Println("Found number:", matches)
	// 	// }

	// 	var formattedData = fmt.Sprintf(`
	// device-mqtt-%s:
	// 	command: /device-mqtt -cp=consul.http://edgex-core-consul:8500 --registry --confdir=/res
	// 	container_name: edgex-device-mqtt-%s
	// 	depends_on:
	// 	- consul
	// 	- data
	// 	- metadata
	// 	- security-bootstrapper
	// 	- mqtt-broker
	// 	entrypoint:
	// 	- /edgex-init/ready_to_run_wait_install.sh
	// 	environment:
	// 		API_GATEWAY_HOST: edgex-kong
	// 		API_GATEWAY_STATUS_PORT: '8100'
	// 		CLIENTS_CORE_COMMAND_HOST: edgex-core-command
	// 		CLIENTS_CORE_DATA_HOST: edgex-core-data
	// 		CLIENTS_CORE_METADATA_HOST: edgex-core-metadata
	// 		CLIENTS_SUPPORT_NOTIFICATIONS_HOST: edgex-support-notifications
	// 		CLIENTS_SUPPORT_SCHEDULER_HOST: edgex-support-scheduler
	// 		DATABASES_PRIMARY_HOST: edgex-redis
	// 		EDGEX_SECURITY_SECRET_STORE: "true"
	// 		MESSAGEQUEUE_HOST: edgex-redis
	// 		PROXY_SETUP_HOST: edgex-security-proxy-setup
	// 		REGISTRY_HOST: edgex-core-consul
	// 		SECRETSTORE_HOST: edgex-vault
	// 		SECRETSTORE_PORT: '8200'
	// 		SERVICE_HOST: edgex-device-mqtt-%s
	// 		SPIFFE_ENDPOINTSOCKET: /tmp/edgex/secrets/spiffe/public/api.sock
	// 		SPIFFE_TRUSTBUNDLE_PATH: /tmp/edgex/secrets/spiffe/trust/bundle
	// 		SPIFFE_TRUSTDOMAIN: edgexfoundry.org
	// 		STAGEGATE_BOOTSTRAPPER_HOST: edgex-security-bootstrapper
	// 		STAGEGATE_BOOTSTRAPPER_STARTPORT: '54321'
	// 		STAGEGATE_DATABASE_HOST: edgex-redis
	// 		STAGEGATE_DATABASE_PORT: '6379'
	// 		STAGEGATE_DATABASE_READYPORT: '6379'
	// 		STAGEGATE_KONGDB_HOST: edgex-kong-db
	// 		STAGEGATE_KONGDB_PORT: '5432'
	// 		STAGEGATE_KONGDB_READYPORT: '54325'
	// 		STAGEGATE_READY_TORUNPORT: '54329'
	// 		STAGEGATE_REGISTRY_HOST: edgex-core-consul
	// 		STAGEGATE_REGISTRY_PORT: '8500'
	// 		STAGEGATE_REGISTRY_READYPORT: '54324'
	// 		STAGEGATE_SECRETSTORESETUP_HOST: edgex-security-secretstore-setup
	// 		STAGEGATE_SECRETSTORESETUP_TOKENS_READYPORT: '54322'
	// 		STAGEGATE_WAITFOR_TIMEOUT: 60s
	// 		MQTTBROKERINFO_HOST: %s
	// 		DEVICE_DEVICESDIR: /res/devices
	// 		DEVICE_PROFILESDIR: /res/profiles
	// 	hostname: edgex-device-mqtt-%s
	// 	image: edgexfoundry/device-mqtt:2.3.0
	// 	networks:
	// 		edgex-network: {}
	// 	ports:
	// 	- 127.0.0.1:59900:59900/tcp
	// 	read_only: true
	// 	restart: always
	// 	security_opt:
	// 	- no-new-privileges:true
	// 	user: 2002:2001
	// 	volumes:
	// 	- edgex-init:/edgex-init:ro,z
	// 	- /tmp/edgex/secrets/device-mqtt:/tmp/edgex/secrets/device-mqtt:ro,z
	// 	- ./device-mqtt-go/cmd/res:/res
	// 	`, deviceInfo.Broker, deviceInfo.Broker, deviceInfo.Broker, deviceInfo.Broker, deviceInfo.Broker)

	// 	if _, err := file.Write([]byte(formattedData)); err != nil {
	// 		WriteError(w, http.StatusBadRequest, err)
	// 		return
	// 	}
	// }

	// url_edgex := fmt.Sprintf("http://localhost:59881/deviceservice/name/device-mqtt", deviceInfo.Broker)
	url_edgex := fmt.Sprintf("http://localhost:59881/api/v2/deviceservice/name/device-mqtt")

	send := false
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	// fmt.Println("Here")
	for i := 0; i < 3; i++ {
		// Send the GET request
		resp, err := c.Get(url_edgex)
		if err != nil {
			// Handle the error (e.g., network issue, service down)
			fmt.Println("Waiting for service to begin...")
		} else {
			// Check if the response status code is 200
			if resp.StatusCode == 200 {
				send = true
				// fmt.Println("Send success")
				resp.Body.Close() // Close the response body when done
				break
			} else {
				// If status code is not 200, handle the error or retry
				fmt.Printf("Received status %d. Retrying...\n", resp.StatusCode)
				resp.Body.Close() // Close the response body even on failure
			}
		}
	}
	if !send {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("Service is not available"))
		return
	}

	var device_name []string
	var profile_name []string

	// fmt.Println(len(deviceInfo.Topic))

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

	url_profile := fmt.Sprintf("http://localhost:59881/api/v2/deviceprofile")

	var profiles []map[string]interface{}

	// Iterate over the device names to create a profile for each one
	// Iterate over the device names to create a profile for each one
	for i := 0; i < len(device_name); i++ {
		var url_resource = fmt.Sprintf("http://localhost:59881/api/v2/deviceresource/profile/%s/resource/%s", fmt.Sprintf("%s-MQTT-device-profile", device_name[i]))

		resp, err := c.Get(url_resource)
		if err != nil {
			// Handle the error (e.g., network issue, service down)
			fmt.Println("Waiting for service to begin...")
		} else {
			// Check if the response status code is 200
			if resp.StatusCode == 200 {
				send = true
				// fmt.Println("Send success")
				resp.Body.Close() // Close the response body when done
				break
			} else {
				// If status code is not 200, handle the error or retry
				fmt.Printf("Received status %d. Retrying...\n", resp.StatusCode)
				resp.Body.Close() // Close the response body even on failure
			}
		}

		data := map[string]interface{}{
			"apiVersion": "v2",
			"profile": map[string]interface{}{
				"name":         fmt.Sprintf("%s-MQTT-device-profile", device_name[i]), // Dynamically set the profile name
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
		}

		// Append the profile data to the profiles slice
		profiles = append(profiles, data)
	}

	// Marshal the map to JSON to verify the structure
	jsonDataProfile, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	reader := bytes.NewReader(jsonDataProfile)
	resp, err := c.Post(url_profile, "application/json", reader)
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("Error add device"))
		fmt.Println(err)
		return
	}
	if resp.StatusCode/100 != 2 { // 200 OK
		WriteError(w, resp.StatusCode, fmt.Errorf("Received non-OK status code: %d", resp.StatusCode))
		fmt.Println(err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error reading response body"))
		fmt.Println(err)
		return
	}

	// Parse the JSON response body into a slice of maps (assuming the response is an array of objects)
	var responseData []map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error parsing response body"))
		fmt.Println(err)
		return
	}

	for _, item := range responseData {
		// Type assertion to extract the statusCode
		if statusCode, ok := item["statusCode"].(float64); ok { // Use float64 because JSON numbers are usually parsed as float64
			if statusCode != http.StatusOK {
				WriteError(w, http.StatusBadRequest, fmt.Errorf("Error adding profile device"))
				return
			}
			fmt.Printf("StatusCode: %.0f\n", statusCode)
		} else {
			WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error parsing response body"))
			fmt.Println("StatusCode not found in map")
			return
		}
	}

	// Print the parsed response (or handle it according to your needs)
	// fmt.Println("Response data:", responseData)

	url_device := fmt.Sprintf("http://localhost:59881/api/v2/device")

	var devices []map[string]interface{}

	// fmt.Println(len(device_name))

	// Loop through the arrays and construct the JSON-like structure
	for i := 0; i < len(device_name); i++ {
		device := map[string]interface{}{
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
				"profileName": fmt.Sprintf("%s-MQTT-device-profile", device_name[i]),
				"protocols": map[string]interface{}{
					"mqtt": map[string]string{
						"CommandTopic": deviceInfo.Topic[i],
					},
				},
			},
		}
		// Add the device map to the devices slice
		devices = append(devices, device)
	}

	// Marshal the final result into JSON
	jsonDataDevice, err := json.MarshalIndent(devices, "", "  ")

	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	reader = bytes.NewReader(jsonDataDevice)
	resp, err = c.Post(url_device, "application/json", reader)
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("Error add device"))
		return
	}
	if resp.StatusCode != http.StatusOK { // 200 OK
		WriteError(w, resp.StatusCode, fmt.Errorf("Received non-OK status code: %d", resp.StatusCode))
		return
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error reading response body"))
		fmt.Println(err)
		return
	}

	// Parse the JSON response body into a slice of maps (assuming the response is an array of objects)
	var responseDataDevice []map[string]interface{}
	err = json.Unmarshal(body, &responseDataDevice)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error parsing response body"))
		fmt.Println(err)
		return
	}

	for _, item := range responseDataDevice {
		// Type assertion to extract the statusCode
		if statusCode, ok := item["statusCode"].(float64); ok { // Use float64 because JSON numbers are usually parsed as float64
			if statusCode != http.StatusOK {
				WriteError(w, http.StatusBadRequest, fmt.Errorf("Error adding device"))
				return
			}
			fmt.Printf("StatusCode: %.0f\n", statusCode)
		} else {
			WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error parsing response body"))
			fmt.Println("StatusCode not found in map")
			return
		}
	}
	WriteJSON(w, http.StatusOK, "Add success")
	// url_profile := fmt.Sprintf("http://localhost:59881/device")
	// device_body :=

	// fmt.Println("Data: ", deviceInfo.Broker)

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
