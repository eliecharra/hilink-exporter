# hilink-exporter

Hilink exporter to expose to prometheus Huawei devices infos

## Usage

```shell
$ docker run eliecharra/hilink-exporter -url="http://192.168.8.1/" -user=admin -password=admin
INFO[0000] Listening on :9770
```

See `--help` for usage details

## Metrics

Here is a **non exhaustive** list of exposed metrics

```
# HELP hilink_device_info device info
# TYPE hilink_device_info gauge
hilink_device_info{DeviceName="B525s-23a",HardwareVersion="WL1B520FM",Imei="XXXXXXXXXXXXXXX",MacAddress1="XX:XX:XX:XX:XX:XX",MacAddress2="",SoftwareVersion="11.189.61.00.1217"} 1
# HELP hilink_traffic_current_download CurrentDownload (bits)
# TYPE hilink_traffic_current_download counter
hilink_traffic_current_download 0
# HELP hilink_traffic_current_download_rate CurrentDownloadRate (bits/s)
# TYPE hilink_traffic_current_download_rate counter
hilink_traffic_current_download_rate 0
# HELP hilink_traffic_current_month_download CurrentMonthDownload (bits)
# TYPE hilink_traffic_current_month_download counter
hilink_traffic_current_month_download{month_last_clear_time="2020-11-30"} 1.054607504e+09
# HELP hilink_traffic_current_month_upload CurrentMonthUpload (bits)
# TYPE hilink_traffic_current_month_upload counter
hilink_traffic_current_month_upload{month_last_clear_time="2020-11-30"} 8.9835791e+07
# HELP hilink_traffic_current_upload CurrentUpload (bits)
# TYPE hilink_traffic_current_upload counter
hilink_traffic_current_upload 0
# HELP hilink_traffic_current_upload_rate CurrentUploadRate (bits/s)
# TYPE hilink_traffic_current_upload_rate counter
hilink_traffic_current_upload_rate 0
# HELP hilink_traffic_total_download TotalDownload (bits)
# TYPE hilink_traffic_total_download counter
hilink_traffic_total_download 1.054607504e+09
# HELP hilink_traffic_total_upload TotalUpload (bits)
# TYPE hilink_traffic_total_upload counter
hilink_traffic_total_upload 8.9835791e+07
# HELP hilink_wan wan info
# TYPE hilink_wan gauge
hilink_wan{WanIPAddress="",WanIPv6Address=""} 1
```

## Tested devices

:warning: list not exhaustive

- B525s-23a