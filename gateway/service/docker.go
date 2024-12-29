package service

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func scanAndUpdate(fpath string, content string) error {
	// Open the file for reading
	file, errFile := os.Open(fpath)
	if errFile != nil {
		return errFile
	}
	defer file.Close()

	// Create a temporary file for writing
	tmpFile, err := os.Create(fpath + ".tmp")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	var foundToken bool
	var targetLine string
	var foundSecret bool
	var foundACL bool
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains ADD_SECRETSTORE_TOKENS
		if strings.Contains(line, "ADD_SECRETSTORE_TOKENS") {
			foundToken = true
			targetLine = line

			// Update the target line
			parts := strings.Split(targetLine, ":")
			if len(parts) > 1 {
				// Trim spaces and add quotes around the value
				secretStoreToken := strings.TrimSpace(parts[1])

				fmt.Println(secretStoreToken)

				re := regexp.MustCompile(`'([^']+)`)

				// Find the match and capture the content inside the quotes
				match := re.FindStringSubmatch(secretStoreToken)

				// If a match is foundToken, the content inside the quotes is at index 1
				if len(match) > 1 {
					fmt.Println(match[1]) // Prints 'abc' without quotes
				} else {
					fmt.Println("No match foundToken")
				}

				// Append additional tokens to the existing value
				secretStoreToken = match[1] + ", " + content

				// Format the line with the updated token values
				updatedLine := "      ADD_SECRETSTORE_TOKENS: '" + secretStoreToken + "'"

				// Write the updated line to the temporary file
				_, err := tmpFile.WriteString(updatedLine + "\n")
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("Docker file has problems")
			}
		} else if strings.Contains(line, "ADD_KNOWN_SECRETS") {
			foundSecret = true
			targetLine = line

			// Update the target line
			parts := strings.Split(targetLine, ":")
			if len(parts) > 1 {
				// Trim spaces and add quotes around the value
				secretStoreToken := strings.TrimSpace(parts[1])

				fmt.Println(secretStoreToken)

				// If a match is foundToken, the content inside the quotes is at index 1

				// Append additional tokens to the existing value
				secretStoreToken = secretStoreToken + ", redisdb[" + content + "]"

				fmt.Sprintln(secretStoreToken)

				// Format the line with the updated token values
				updatedLine := "      ADD_KNOWN_SECRETS: " + secretStoreToken + ""

				// Write the updated line to the temporary file
				_, err := tmpFile.WriteString(updatedLine + "\n")
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("Docker file has problems")
			}
		} else if strings.Contains(line, "ADD_REGISTRY_ACL_ROLES") {
			foundACL = true
			targetLine = line

			// Update the target line
			parts := strings.Split(targetLine, ":")
			if len(parts) > 1 {
				// Trim spaces and add quotes around the value
				secretStoreToken := strings.TrimSpace(parts[1])

				fmt.Println(secretStoreToken)

				re := regexp.MustCompile(`'([^']+)`)

				// Find the match and capture the content inside the quotes
				match := re.FindStringSubmatch(secretStoreToken)

				// If a match is foundToken, the content inside the quotes is at index 1
				if len(match) > 1 {
					fmt.Println(match[1]) // Prints 'abc' without quotes
				} else {
					fmt.Println("No match foundToken")
				}

				// Append additional tokens to the existing value
				secretStoreToken = match[1] + ", " + content

				// Format the line with the updated token values
				updatedLine := "      ADD_REGISTRY_ACL_ROLES: '" + secretStoreToken + "'"

				// Write the updated line to the temporary file
				_, err := tmpFile.WriteString(updatedLine + "\n")
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("Docker file has problems")
			}
		} else {
			// Write the line to the temporary file
			_, err := tmpFile.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if !foundACL {
		return fmt.Errorf("Docker compose does not have ADD_REGISTRY_ACL_ROLES field")
	}
	if !foundToken {
		return fmt.Errorf("Docker compose does not have ADD_SECRETSTORE_TOKENS field")
	}
	if !foundSecret {
		return fmt.Errorf("Docker compose does not have ADD_KNOWN_SECRETS field")
	}

	// Replace the original file with the temporary one
	err = os.Remove(fpath)
	if err != nil {
		return err
	}
	err = os.Rename(fpath+".tmp", fpath)
	if err != nil {
		return err
	}

	return nil
}
