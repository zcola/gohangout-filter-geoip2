# gohangout-filter-geoip2
```
filters:
    - Json:
        field: message
        remove_fields: ['message']
    - '/mnt/c/work/gohangout_1.8.1/gohangout/geoip2_new.so':
        src: '[proxima_meta][ip]'
        dbPath: '/mnt/c/work/gohangout_1.8.1/gohangout/dasds.mmdb'
        target: ip_geo
    - '/mnt/c/work/gohangout_1.8.1/gohangout/geoip2_new.so':
        src: '[proxima_meta][ip]'
        dbPath: '/mnt/c/work/gohangout_1.8.1/gohangout/das_isp.mmdb'
        target: ip_geo
```
