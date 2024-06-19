package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type Status struct {
	statusCode int
	message    string
}

func handleProjectDeploy(c *gin.Context) {
	projectName := c.Param("projectName")

	scriptPath := "./deploy_scripts/" + projectName + ".sh"
	_, err := os.Stat(scriptPath)
	if os.IsNotExist(err) {
		log.Fatalf("Deploy script %s does not exist", scriptPath)
		c.AbortWithStatus(404)
	}

	// Special handling for deployinator, as it needs to respond to the request before running it's own deploy script
	if projectName == "deployinator" {
		go deployProject(scriptPath)
		c.String(200, "Attempting to deploy deployinator")
	} else {
		result := deployProject(scriptPath)
		c.String(result.statusCode, result.message)
	}
}

func deployProject(scriptPath string) Status {
	scriptStatus := make(chan Status)
	go runDeployScript(scriptPath, scriptStatus)
	result := <-scriptStatus
	return result
}

func runDeployScript(scriptPath string, status chan Status) {
	cmd := exec.Command("/bin/bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		if gin.Mode() == gin.ReleaseMode {
			status <- Status{500, "Error running deploy script"}
		} else {
			status <- Status{500, "Error running deploy script: " + err.Error()}
		}
	}
	status <- Status{200, "Deploy script finished"}
}

func validateSecret(c *gin.Context) {
	reqSecret := getSignatureHeader(c)
	if reqSecret == "" {
		log.Print("No signature header found")
		c.AbortWithStatus(401)
		return
	}

	realSecret, exist := os.LookupEnv("SECRET_KEY")
	if !exist {
		log.Print("SECRET_KEY not found in .env file")
		c.AbortWithStatus(500)
	}

	payloadBody, err := getPayloadBody(c)
	if err != nil {
		log.Printf("Error getting payload body: %s", err)
		c.AbortWithStatus(500)
	}

	if err := verifySignature(payloadBody, realSecret, reqSecret); err != nil {
		log.Printf("Error verifying signature: %s", err)
		c.AbortWithStatus(401)
	}
}

// verifySignature verifies that the payload was sent from GitHub by validating the SHA256 signature.
func verifySignature(payloadBody []byte, secretToken string, signatureHeader string) error {

	// Create a new HMAC using the secret token and SHA256 hash
	h := hmac.New(sha256.New, []byte(secretToken))
	h.Write(payloadBody)
	expectedSignature := "sha256=" + hex.EncodeToString(h.Sum(nil))

	// Compare the expected signature with the provided signature header
	if !hmac.Equal([]byte(expectedSignature), []byte(signatureHeader)) {
		return errors.New("request signatures didn't match")
	}

	return nil
}

func getPayloadBody(c *gin.Context) ([]byte, error) {
	payloadBody, err := c.GetRawData()
	if err != nil {
		return nil, err
	}
	return payloadBody, nil
}

func getSignatureHeader(c *gin.Context) string {
	return c.GetHeader("X-Hub-Signature-256")
}
