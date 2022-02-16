package builder

import "testing"

func TestShellCommand(t *testing.T) {
	cmd := DockerShellCommand().
		Container("my-fancy-name").
		Shell("bash").
		Build()

	expected := "docker exec -it my-fancy-name bash"

	if cmd != expected {
		t.Errorf("Expected %s, got %s", expected, cmd)
	}
}

func TestShellCommandOverwrites(t *testing.T) {
	cmd := DockerShellCommand().
		Container("my-fancy-name").
		Shell("aaa").
		Shell("bbb").
		Build()

	expected := "docker exec -it my-fancy-name bbb"

	if cmd != expected {
		t.Errorf("Expected %s, got %s", expected, cmd)
	}
}
