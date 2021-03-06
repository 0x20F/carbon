package builder

import "testing"

func TestComposeCommand(t *testing.T) {
	cmd := DockerComposeCommand().
		File("docker-compose.yml").
		Up().
		Build()

	expected := "docker compose -f docker-compose.yml up"

	if cmd != expected {
		t.Errorf("Expected %s, got %s", expected, cmd)
	}
}

func TestComposeCommandUniqueValues(t *testing.T) {
	cmd := DockerComposeCommand().
		File("docker-compose.yml").
		Up().
		Down(). // We want the latest unique to override all others with the same priority
		Build()

	expected := "docker compose -f docker-compose.yml down"

	if cmd != expected {
		t.Errorf("Expected %s, got %s", expected, cmd)
	}
}

func TestComposeCommandWithSingleService(t *testing.T) {
	cmd := DockerComposeCommand().
		File("docker-compose.yml").
		Service("web").
		Up().
		Build()

	expected := "docker compose -f docker-compose.yml up web"

	if cmd != expected {
		t.Errorf("Expected %s, got %s", expected, cmd)
	}
}

func TestComposeCommandWithMultipleServices(t *testing.T) {
	cmd := DockerComposeCommand().
		File("docker-compose.yml").
		Service("web").
		Service("db").
		Restart().
		Build()

	expected := "docker compose -f docker-compose.yml restart web db"

	if cmd != expected {
		t.Errorf("Expected %s, got %s", expected, cmd)
	}
}
