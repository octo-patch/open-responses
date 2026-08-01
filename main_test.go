package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMinimaxEnvVars(t *testing.T) {
	values := make(map[string]string)
	for _, env := range minimaxEnvVars() {
		values[env.Name] = env.Default
	}

	if values["MINIMAX_API_KEY"] != "" {
		t.Fatalf("MINIMAX_API_KEY default = %q, want empty", values["MINIMAX_API_KEY"])
	}
	if values["MINIMAX_BASE_URL"] != minimaxGlobalBaseURL {
		t.Fatalf("MINIMAX_BASE_URL default = %q, want %q", values["MINIMAX_BASE_URL"], minimaxGlobalBaseURL)
	}
	if values["MINIMAX_MODEL_ID"] != minimaxDefaultModel {
		t.Fatalf("MINIMAX_MODEL_ID default = %q, want %q", values["MINIMAX_MODEL_ID"], minimaxDefaultModel)
	}
}

func TestCreateEnvFileIncludesMinimaxProviderConfig(t *testing.T) {
	originalConfig := config
	t.Cleanup(func() { config = originalConfig })

	envPath := filepath.Join(t.TempDir(), ".env")
	config = Config{
		Host:       "127.0.0.1",
		Port:       "8080",
		DockerTag:  "test",
		APIVersion: "test",
		EnvFile:    envPath,
	}

	createEnvFile(map[string]string{
		"MINIMAX_API_KEY":  "",
		"MINIMAX_BASE_URL": minimaxGlobalBaseURL,
		"MINIMAX_MODEL_ID": minimaxDefaultModel,
	})

	data, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("read generated env file: %v", err)
	}
	content := string(data)

	expected := []string{
		"# Models: " + minimaxDefaultModel + ", " + minimaxFallbackModel,
		"# Global endpoint: " + minimaxGlobalBaseURL,
		"# China endpoint: " + minimaxChinaBaseURL,
		`MINIMAX_API_KEY=""`,
		`MINIMAX_BASE_URL="` + minimaxGlobalBaseURL + `"`,
		`MINIMAX_MODEL_ID="` + minimaxDefaultModel + `"`,
	}
	for _, value := range expected {
		if !strings.Contains(content, value) {
			t.Errorf("generated env file does not contain %q", value)
		}
	}
}

func TestDefaultComposeTemplateIncludesMinimaxProviderConfig(t *testing.T) {
	template := getDefaultComposeTemplate()
	expected := []string{
		"MINIMAX_API_KEY: ${MINIMAX_API_KEY}",
		"MINIMAX_BASE_URL: ${MINIMAX_BASE_URL:-" + minimaxGlobalBaseURL + "}",
		"MINIMAX_MODEL_ID: ${MINIMAX_MODEL_ID:-" + minimaxDefaultModel + "}",
	}
	for _, value := range expected {
		if !strings.Contains(template, value) {
			t.Errorf("default Compose template does not contain %q", value)
		}
	}
}
