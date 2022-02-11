package builder

import "testing"

func TestBuildCommand(t *testing.T) {
	cmd := DockerBuildCommand().
		Tag("thing:latest").
		File("Dockerfile").
		Path(".").
		BuildArg("GITHUB_TOKEN=1241234").
		Build()

	expected := "docker build -f Dockerfile -t thing:latest --build-arg GITHUB_TOKEN=1241234 ."

	if cmd != expected {
		t.Errorf("Expected %s, got %s", expected, cmd)
	}
}
