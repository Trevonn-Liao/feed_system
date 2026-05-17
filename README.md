# Feed System

## Configuration

The project supports YAML configuration files and environment variable overrides.

Configuration priority:

1. Environment variables
2. `configs/dev_local.yaml`
3. `configs/config.local.yaml`
4. `configs/config.example.yaml`
5. Built-in development defaults

Create a local secret config file:

```bash
cp configs/config.example.yaml configs/dev_local.yaml
```

Generate JWT secrets:

```bash
openssl rand -base64 32
openssl rand -base64 32
```

Put the generated values into `configs/dev_local.yaml`:

```yaml
jwt:
  access_secret: "your-generated-access-secret"
  refresh_secret: "your-generated-refresh-secret"
```

You can also override JWT secrets with environment variables:

```bash
JWT_ACCESS_SECRET=your-generated-access-secret
JWT_REFRESH_SECRET=your-generated-refresh-secret
```

Do not commit local secrets. `configs/dev_local.yaml`, `configs/config.local.yaml`, and `.env` must never be pushed to GitHub.

## RedisBloom

The user Bloom filter uses the RedisBloom module commands:

- `BF.RESERVE`
- `BF.ADD`
- `BF.EXISTS`

Use Redis Stack or install the RedisBloom module before starting the service. A plain Redis server without RedisBloom will fail during startup.
