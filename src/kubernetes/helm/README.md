# CinemaAbyss Helm Chart

This Helm chart deploys the CinemaAbyss application on a Kubernetes cluster.

## Prerequisites

- Kubernetes 1.16+
- Helm 3.0+
- PV provisioner support in the underlying infrastructure (if persistence is enabled)

## Installing the Chart

To install the chart with the release name `cinemaabyss`:

```bash
helm install cinemaabyss ./cinemaabyss
```

The command deploys CinemaAbyss on the Kubernetes cluster with default configuration. The [Parameters](#parameters) section lists the parameters that can be configured during installation.

## Uninstalling the Chart

To uninstall/delete the `cinemaabyss` deployment:

```bash
helm uninstall cinemaabyss
```

## Parameters

### Global Parameters

| Name                | Description                                     | Value           |
|---------------------|-------------------------------------------------|-----------------|
| `global.namespace`  | Namespace to deploy all resources               | `cinemaabyss`   |
| `global.domain`     | Domain name for the application                 | `cinemaabyss.example.com` |

### Database Parameters

| Name                           | Description                                     | Value           |
|--------------------------------|-------------------------------------------------|-----------------|
| `database.host`                | PostgreSQL host                                 | `postgres`      |
| `database.port`                | PostgreSQL port                                 | `5432`          |
| `database.name`                | PostgreSQL database name                        | `cinemaabyss`   |
| `database.user`                | PostgreSQL username                             | `postgres`      |
| `database.password`            | PostgreSQL password (base64 encoded)            | `cG9zdGdyZXNfcGFzc3dvcmQ=` |
| `database.image.repository`    | PostgreSQL image repository                     | `postgres`      |
| `database.image.tag`           | PostgreSQL image tag                            | `14`            |
| `database.image.pullPolicy`    | PostgreSQL image pull policy                    | `IfNotPresent`  |
| `database.resources.limits.cpu`| PostgreSQL CPU limit                            | `1000m`         |
| `database.resources.limits.memory` | PostgreSQL memory limit                     | `1Gi`           |
| `database.resources.requests.cpu` | PostgreSQL CPU request                       | `500m`          |
| `database.resources.requests.memory` | PostgreSQL memory request                 | `512Mi`         |
| `database.persistence.enabled` | Enable persistence for PostgreSQL               | `true`          |
| `database.persistence.size`    | PostgreSQL PVC size                             | `10Gi`          |
| `database.persistence.accessMode` | PostgreSQL PVC access mode                   | `ReadWriteOnce` |

### Monolith Parameters

| Name                           | Description                                     | Value           |
|--------------------------------|-------------------------------------------------|-----------------|
| `monolith.enabled`             | Enable monolith deployment                      | `true`          |
| `monolith.image.repository`    | Monolith image repository                       | `ghcr.io/db-exp/cinemaabysstest/monolith` |
| `monolith.image.tag`           | Monolith image tag                              | `latest`        |
| `monolith.image.pullPolicy`    | Monolith image pull policy                      | `Always`        |
| `monolith.replicas`            | Number of monolith replicas                     | `1`             |
| `monolith.resources.limits.cpu`| Monolith CPU limit                              | `500m`          |
| `monolith.resources.limits.memory` | Monolith memory limit                       | `512Mi`         |
| `monolith.resources.requests.cpu` | Monolith CPU request                         | `100m`          |
| `monolith.resources.requests.memory` | Monolith memory request                   | `128Mi`         |
| `monolith.service.port`        | Monolith service port                           | `8080`          |
| `monolith.service.targetPort`  | Monolith container port                         | `8080`          |
| `monolith.service.type`        | Monolith service type                           | `ClusterIP`     |

### Proxy Service Parameters

| Name                           | Description                                     | Value           |
|--------------------------------|-------------------------------------------------|-----------------|
| `proxyService.enabled`         | Enable proxy service deployment                 | `true`          |
| `proxyService.image.repository`| Proxy service image repository                  | `ghcr.io/db-exp/cinemaabysstest/proxy-service` |
| `proxyService.image.tag`       | Proxy service image tag                         | `latest`        |
| `proxyService.image.pullPolicy`| Proxy service image pull policy                 | `Always`        |
| `proxyService.replicas`        | Number of proxy service replicas                | `1`             |
| `proxyService.resources.limits.cpu`| Proxy service CPU limit                     | `300m`          |
| `proxyService.resources.limits.memory` | Proxy service memory limit              | `256Mi`         |
| `proxyService.resources.requests.cpu` | Proxy service CPU request                | `100m`          |
| `proxyService.resources.requests.memory` | Proxy service memory request          | `128Mi`         |
| `proxyService.service.port`    | Proxy service port                              | `80`            |
| `proxyService.service.targetPort` | Proxy service container port                 | `8000`          |
| `proxyService.service.type`    | Proxy service type                              | `ClusterIP`     |

### Movies Service Parameters

| Name                           | Description                                     | Value           |
|--------------------------------|-------------------------------------------------|-----------------|
| `moviesService.enabled`        | Enable movies service deployment                | `true`          |
| `moviesService.image.repository`| Movies service image repository                | `ghcr.io/db-exp/cinemaabysstest/movies-service` |
| `moviesService.image.tag`      | Movies service image tag                        | `latest`        |
| `moviesService.image.pullPolicy`| Movies service image pull policy               | `Always`        |
| `moviesService.replicas`       | Number of movies service replicas               | `1`             |
| `moviesService.resources.limits.cpu`| Movies service CPU limit                   | `300m`          |
| `moviesService.resources.limits.memory` | Movies service memory limit            | `256Mi`         |
| `moviesService.resources.requests.cpu` | Movies service CPU request              | `100m`          |
| `moviesService.resources.requests.memory` | Movies service memory request        | `128Mi`         |
| `moviesService.service.port`   | Movies service port                             | `8081`          |
| `moviesService.service.targetPort` | Movies service container port               | `8081`          |
| `moviesService.service.type`   | Movies service type                             | `ClusterIP`     |

### Events Service Parameters

| Name                           | Description                                     | Value           |
|--------------------------------|-------------------------------------------------|-----------------|
| `eventsService.enabled`        | Enable events service deployment                | `true`          |
| `eventsService.image.repository`| Events service image repository                | `ghcr.io/db-exp/cinemaabysstest/events-service` |
| `eventsService.image.tag`      | Events service image tag                        | `latest`        |
| `eventsService.image.pullPolicy`| Events service image pull policy               | `Always`        |
| `eventsService.replicas`       | Number of events service replicas               | `1`             |
| `eventsService.resources.limits.cpu`| Events service CPU limit                   | `300m`          |
| `eventsService.resources.limits.memory` | Events service memory limit            | `256Mi`         |
| `eventsService.resources.requests.cpu` | Events service CPU request              | `100m`          |
| `eventsService.resources.requests.memory` | Events service memory request        | `128Mi`         |
| `eventsService.service.port`   | Events service port                             | `8082`          |
| `eventsService.service.targetPort` | Events service container port               | `8082`          |
| `eventsService.service.type`   | Events service type                             | `ClusterIP`     |

### Kafka Parameters

| Name                           | Description                                     | Value           |
|--------------------------------|-------------------------------------------------|-----------------|
| `kafka.enabled`                | Enable Kafka deployment                         | `true`          |
| `kafka.image.repository`       | Kafka image repository                          | `wurstmeister/kafka` |
| `kafka.image.tag`              | Kafka image tag                                 | `2.13-2.7.0`    |
| `kafka.image.pullPolicy`       | Kafka image pull policy                         | `IfNotPresent`  |
| `kafka.replicas`               | Number of Kafka replicas                        | `1`             |
| `kafka.resources.limits.cpu`   | Kafka CPU limit                                 | `1000m`         |
| `kafka.resources.limits.memory`| Kafka memory limit                              | `1Gi`           |
| `kafka.resources.requests.cpu` | Kafka CPU request                               | `200m`          |
| `kafka.resources.requests.memory` | Kafka memory request                         | `512Mi`         |
| `kafka.persistence.enabled`    | Enable persistence for Kafka                    | `true`          |
| `kafka.persistence.size`       | Kafka PVC size                                  | `5Gi`           |
| `kafka.persistence.accessMode` | Kafka PVC access mode                           | `ReadWriteOnce` |
| `kafka.topics`                 | Kafka topics configuration                      | See values.yaml |

### Zookeeper Parameters

| Name                           | Description                                     | Value           |
|--------------------------------|-------------------------------------------------|-----------------|
| `zookeeper.enabled`            | Enable Zookeeper deployment                     | `true`          |
| `zookeeper.image.repository`   | Zookeeper image repository                      | `wurstmeister/zookeeper` |
| `zookeeper.image.tag`          | Zookeeper image tag                             | `latest`        |
| `zookeeper.image.pullPolicy`   | Zookeeper image pull policy                     | `IfNotPresent`  |
| `zookeeper.replicas`           | Number of Zookeeper replicas                    | `1`             |
| `zookeeper.resources.limits.cpu`| Zookeeper CPU limit                            | `500m`          |
| `zookeeper.resources.limits.memory` | Zookeeper memory limit                     | `512Mi`         |
| `zookeeper.resources.requests.cpu` | Zookeeper CPU request                       | `100m`          |
| `zookeeper.resources.requests.memory` | Zookeeper memory request                 | `256Mi`         |
| `zookeeper.persistence.enabled`| Enable persistence for Zookeeper                | `true`          |
| `zookeeper.persistence.size`   | Zookeeper PVC size                              | `1Gi`           |
| `zookeeper.persistence.accessMode` | Zookeeper PVC access mode                   | `ReadWriteOnce` |

### Ingress Parameters

| Name                           | Description                                     | Value           |
|--------------------------------|-------------------------------------------------|-----------------|
| `ingress.enabled`              | Enable ingress                                  | `true`          |
| `ingress.className`            | Ingress class name                              | `nginx`         |
| `ingress.annotations`          | Ingress annotations                             | See values.yaml |
| `ingress.hosts`                | Ingress hosts configuration                     | See values.yaml |

### Application Configuration

| Name                           | Description                                     | Value           |
|--------------------------------|-------------------------------------------------|-----------------|
| `config.gradualMigration`      | Enable gradual migration                        | `true`          |
| `config.moviesMigrationPercent`| Movies migration percentage                     | `100`           |

## Architecture

The CinemaAbyss application consists of the following components:

1. **Monolith**: The main application that handles user authentication, subscriptions, and payments.
2. **Proxy Service**: A service that routes requests to the appropriate microservice or the monolith.
3. **Movies Service**: A microservice that handles movie-related functionality.
4. **Events Service**: A microservice that handles event processing using Kafka.
5. **PostgreSQL**: The database used by all services.
6. **Kafka**: Message broker for event-driven communication.
7. **Zookeeper**: Required for Kafka coordination.

## Persistence

The chart mounts a Persistent Volume for PostgreSQL, Kafka, and Zookeeper. The volume is created using dynamic volume provisioning. If you want to disable this functionality, you can set `database.persistence.enabled`, `kafka.persistence.enabled`, and `zookeeper.persistence.enabled` to `false`.

## Image Pull Secrets

The chart includes a secret for pulling images from private registries. The secret is created using the value provided in `imagePullSecrets.dockerconfigjson`.