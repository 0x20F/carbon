package builder

import "testing"

func TestShellCommand(t *testing.T) {
	cmd := DockerShellCommand().
		Container("my-fancy-name").
		Bash().
		Build()

	expected := "docker exec -it my-fancy-name bash"

	if cmd != expected {
		t.Errorf("Expected %s, got %s", expected, cmd)
	}
}

func TestShellCommandOverwrites(t *testing.T) {
	cmd := DockerShellCommand().
		Container("my-fancy-name").
		Bash().
		Sh().
		Build()

	expected := "docker exec -it my-fancy-name sh"

	if cmd != expected {
		t.Errorf("Expected %s, got %s", expected, cmd)
	}
}

func TestShellCommandCustomShell(t *testing.T) {
	cmd := DockerShellCommand().
		Container("my-fancy-name").
		Shell("/bin/zsh").
		Build()

	expected := "docker exec -it my-fancy-name /bin/zsh"

	if cmd != expected {
		t.Errorf("Expected %s, got %s", expected, cmd)
	}
}
