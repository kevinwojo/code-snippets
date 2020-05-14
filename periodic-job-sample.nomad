job "poke-circle-status" {
  datacenters = ["dc1"]
  type = "batch"
  periodic {
    cron  = "*/30 * * * * *"
    prohibit_overlap = true
  }

  group "poker" {
    count = 1
    task "curl" {
      driver = "docker"
      config {
        image = "curlimages/curl:latest"
        args = ["-I", "https://iscircleciup.net"]
      }

      resources {
        cpu    = 500 # 500 MHz
        memory = 64

        network {
          mbits = 10
        }
      }
    }
  }
}
