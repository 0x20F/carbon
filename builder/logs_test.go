package builder

import "testing"

func TestLogsCommandFollow(t *testing.T) {
	cmd := DockerLogsCommand().
		Follow().
		Container("thing").
		Build()

	expected := "docker logs -f thing"

	if cmd != expected {
		t.Errorf("Expected %s, got %s", expected, cmd)
	}
}

func TestLogsCommandNoFollow(t *testing.T) {
	cmd := DockerLogsCommand().
		Container("thing").
		Build()

	expected := "docker logs thing"

	if cmd != expected {
		t.Errorf("Expected %s, got %s", expected, cmd)
	}
}
