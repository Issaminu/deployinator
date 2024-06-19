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

func handleProjectDeploy(c *gin.Context) {
	projectName := c.Param("projectName")

	scriptPath := "./deploy_scripts/" + projectName + ".sh"
	_, err := os.Stat(scriptPath)
	if os.IsNotExist(err) {
		log.Fatalf("Deploy script %s does not exist", scriptPath)
	}

	if projectName == "deployinator" {
		go runDeployScript(c, scriptPath)
		c.String(200, "Deployinator is attempting to deploy itself")
		return
	}
	runDeployScript(c, scriptPath)
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

func runDeployScript(c *gin.Context, scriptPath string) {
	cmd := exec.Command("/bin/bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		c.String(500, "Error deploying project")
	}
	c.String(200, "Deployed successfully")
}
