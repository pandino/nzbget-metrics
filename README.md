# nzbget-metrics

A Prometheus exporter for [nzbget](https://nzbget.com/), shipped as a
[linuxserver.io docker mod](https://www.linuxserver.io/blog/2019-09-14-customizing-our-containers).
It runs inside the `linuxserver/nzbget` container as an s6 service and scrapes the local JSON-RPC API.

## Usage

Add the mod to your nzbget container:

```yaml
environment:
  DOCKER_MODS: ghcr.io/youruser/nzbget-metrics:latest
ports:
  - "9452:9452"
```

Zero additional configuration is needed — the exporter reads `/config/nzbget.conf` automatically.

## Configuration

All settings have sensible defaults derived from nzbget's config file. Override with environment variables:

| Variable          | Default                    | Description                       |
|-------------------|----------------------------|-----------------------------------|
| `NZBGET_CONFIG`   | `/config/nzbget.conf`      | Path to nzbget config file        |
| `NZBGET_HOST`     | `127.0.0.1` (from conf)    | nzbget host                       |
| `NZBGET_PORT`     | `6789` (from `ControlPort`)| nzbget JSON-RPC port              |
| `NZBGET_USERNAME` | `nzbget` (from conf)       | API username                      |
| `NZBGET_PASSWORD` | `tegbzn6789` (from conf)   | API password                      |
| `METRICS_PORT`    | `9452`                     | Port for the `/metrics` endpoint  |
| `METRICS_PATH`    | `/metrics`                 | Path for the metrics endpoint     |

## Metrics

| Metric | Description |
|--------|-------------|
| `nzbget_up` | 1 if the API responded this scrape |
| `nzbget_info{version}` | Version label (constant 1) |
| `nzbget_download_rate_bytes_per_second` | Current download rate |
| `nzbget_average_download_rate_bytes_per_second` | Average download rate |
| `nzbget_download_limit_bytes_per_second` | Speed limit |
| `nzbget_remaining_size_bytes` | Remaining queue size |
| `nzbget_downloaded_size_bytes` | Total downloaded |
| `nzbget_article_cache_bytes` | Article cache usage |
| `nzbget_free_disk_space_bytes` | Free disk space |
| `nzbget_total_disk_space_bytes` | Total disk space |
| `nzbget_uptime_seconds` | Server uptime |
| `nzbget_download_time_seconds` | Total download time |
| `nzbget_threads` | Active thread count |
| `nzbget_post_jobs` | Post-processing jobs |
| `nzbget_urls` | URL fetches in queue |
| `nzbget_queue_scripts` | Queued scripts |
| `nzbget_queued` | NZBs in queue |
| `nzbget_paused{kind}` | 1 if paused (kind: download/server/post/scan) |
| `nzbget_server_standby` | 1 if in standby |
| `nzbget_quota_reached` | 1 if quota reached |
| `nzbget_day_size_bytes` | Downloaded today |
| `nzbget_month_size_bytes` | Downloaded this month |
| `nzbget_news_server_active{id}` | 1 if news server is active |

## Prometheus scrape config

```yaml
scrape_configs:
  - job_name: nzbget
    static_configs:
      - targets: ["nzbget-host:9452"]
```

## Building

```sh
docker build -t nzbget-metrics .
```

Requires Docker Hub secrets `DOCKERHUB_TOKEN` and variable `DOCKERHUB_USERNAME` for CI push.
