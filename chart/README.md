# Chart

## Required additional variables

Please create a `secrets.yaml` file containing the following secrets:

```yaml
jwt_secret: ""
google:
  client_id: ""
  client_secret: ""
  default_admin_id: ""
```

Then, just run `helm install til . --values secrets.yaml`