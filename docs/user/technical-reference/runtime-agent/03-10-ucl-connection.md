# Connection with UCL

Runtime Agent connects to UCL using a one-time token from the Connector and exchanges it for a certificate, which is later used to fetch Applications from the UCL Director.

The initial connection requires the following parameters:

| **Parameter** | **Description** | **Example value** |
|---------------|-----------------|-------------------|
| **CONNECTOR_URL** | Connector URL | `https://compass-gateway.kyma.local/connector/graphql` |
| **RUNTIME_ID** | ID of the Runtime registered in the Director | `1ae04041-17e5-478f-91f8-3a2ddc7700de` |
| **TENANT** | Tenant ID  | `3e64ebae-38b5-46a0-b1ed-9ccee153a0ae` |
| **TOKEN** | One-time token generated for the Runtime | `2I7VVX5CqxHioEBQGPxWSp3k90uw51tmx5dbo0IZd5VNFzGoPfppYrMIuoCNwFOKp05wsioJNLJYxdI-LKlUYA==` |

Runtime Agent reads this configuration from the Secret specified in the Runtime Agent Deployment (`compass-agent-configuration` by default).

To see how to create the Secret, see the [tutorial](./tutorials/01-90-configure-runtime-agent-with-compass.md).

## Connection Status

The connection status is preserved in the [Connection](../../resources/04-20-compassconnection.md) custom resource (CR). This CR also stores the Connector URL and the Director URL.

> [!NOTE]
> Mind that UCL was previously called Compass. For historical reasons, the Connection CR is still called `CompassConnection`.

## Reconnecting

If the connection with UCL fails, Runtime Agent keeps trying to connect with the token from the Secret. If the connection is established successfully, Runtime Agent ignores the Secret until the connection is lost.

To see how to reconnect Runtime Agent with UCL, see this [tutorial](./tutorials/01-100-reconnect-runtime-agent-with-compass.md).
