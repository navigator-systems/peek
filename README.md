# Peek — Monitor del Sistema Linux desde la Terminal

## Tabla de Contenidos

1. [Descripción](#descripción)
2. [Requisitos](#requisitos)
3. [Instalación](#instalación)
4. [Uso](#uso)
   - [peek cpu](#peek-cpu)
   - [peek mem](#peek-mem)
   - [peek stg](#peek-stg)
   - [peek port](#peek-port)
5. [Arquitectura del Proyecto](#arquitectura-del-proyecto)
   - [Estructura de Directorios](#estructura-de-directorios)
   - [Flujo de Ejecución](#flujo-de-ejecución)
6. [Dependencias](#dependencias)
7. [Guía de Desarrollo](#guía-de-desarrollo)
   - [Añadir un Nuevo Comando](#añadir-un-nuevo-comando)
   - [Convenciones de Código](#convenciones-de-código)
8. [Ejemplos de Salida](#ejemplos-de-salida)
9. [Limitaciones](#limitaciones)
10. [Licencia](#licencia)

---

## Descripción

**Peek** es una herramienta de línea de comandos escrita en **Go** que proporciona información en tiempo real sobre recursos del sistema operativo Linux. Diseñada para ser minimalista, rápida y fácil de usar, permite consultar el estado de CPU, memoria, almacenamiento y puertos de red sin necesidad de herramientas externas complejas.

```
peek cpu       # Información del CPU y sistema operativo
peek mem       # Estadísticas de memoria virtual
peek stg       # Uso de almacenamiento por partición
peek port      # Listar puertos en escucha
peek port 8080 # Filtrar por puerto específico
```

---

## Requisitos

| Requisito | Versión Mínima | Notas |
|-----------|---------------|-------|
| **Go** | 1.27 | Requerido para compilar el proyecto |

---

## Instalación

### 1. Clonar el Repositorio

```bash
git clone https://github.com/navigator/peek.git
cd peek
```

### 2. Descargar Dependencias

```bash
go mod tidy
```

### 3. Compilar

```bash
# Compilación estándar
go build -o peek .

# Compilación conflags de optimización
go build -ldflags="-s -w" -o peek .
```

### 4. Ejecutar

```bash
# Ejecutar directamente con Go
go run .

# O usar el binario compilado
./peek --help
```

### 5. Instalación Global (Opcional)

```bash
# Mover el binario a un directorio en PATH
sudo mv peek /usr/local/bin/

# Verificar instalación
peek --help
```

---

## Uso

### `peek cpu`

Muestra información detallada sobre el procesador y el sistema operativo host.

```bash
peek cpu
```

**Salida esperada:**

```
OS: Linux - Platform debian
Kernel: 6.1.0 - x86_64
CPU: Intel(R) Core(TM) i7-10700K CPU @ 3.80GHz - Cores 16
```

| Campo | Descripción |
|-------|-------------|
| OS | Sistema operativo base |
| Platform | Familia/distribución del sistema |
| Kernel | Versión del kernel de Linux |
| KernelArch | Arquitectura del kernel (x86_64, arm64, etc.) |
| CPU | Modelo completo del procesador |
| Cores | Cantidad de núcleos lógicos disponibles |

---

### `peek mem`

Proporciona estadísticas de memoria virtual del sistema.

```bash
peek mem
```

**Salida esperada:**

```
Total: 31.17 GiB, Free: 8.24 GiB, Used: 22.93 GiB, UsedPercent: 73.561%
```

| Campo | Descripción |
|-------|-------------|
| Total | Memoria total disponible |
| Free | Memoria completamente libre |
| Used | Memoria en uso |
| UsedPercent | Porcentaje de uso |

---

### `peek stg`

Lista todas las particiones montadas con su uso de espacio en disco.

```bash
peek stg
```

**Salida esperada:**

```
MOUNT           TYPE      TOTAL     USED     FREE
/               ext4      50.00 GB  25.10 GB 24.90 GB
/boot/efi       vfat      1.00 GB   15.05 MB 989.60 MB
/home           ext4      200.00 GB 80.00 GB 120.00 GB

TOTAL STORAGE OVERVIEW
Total: 251.00 GB  Used: 105.35 GB  Free: 145.65 GB
```

| Campo | Descripción |
|-------|-------------|
| MOUNT | Punto de montaje de la partición |
| TYPE | Sistema de archivos (ext4, xfs, btrfs, etc.) |
| TOTAL | Espacio total de la partición |
| USED | Espacio utilizado |
| FREE | Espacio disponible |

---

### `peek port`

Lista todos los sockets en estado LISTEN (puertos abiertos esperando conexiones).

```bash
# Listar todos los puertos
peek port

# Filtrar por puerto específico
peek port 8080
peek port 443
```

**Salida esperada:**

```
Local Address        Proto  PID    Process
0.0.0.0:22           TCP    1234   sshd
127.0.0.1:5037       TCP    31850  adb
0.0.0.0:8080         TCP    9012   nginx
:::1716              TCP    2247   kdeconnected
```

| Campo | Descripción |
|-------|-------------|
| Local Address | Dirección IP y puerto (formato IP:Puerto) |
| Proto | Protocolo de red (TCP, UDP, UNIX) |
| PID | Identificador del proceso |
| Process | Nombre del proceso que mantiene el puerto |

---

### `peek proc`

Lista todos los procesos activos.

```bash
# Listar todos los procesos
peek proc

# Filtrar por los 10 procesos que mas consumen CPU
peek proc cpu

# Filtrar por los 10 procesos que mas consumen Memoria
peek proc mem
```

**Salida esperada:**

```
Name                 PID    CPU    Mem    Command             
peek                 15730  46.524 0.070 ./peek proc cpu     
brave                2981   11.246 1.061 /opt/brave.com/brave/brave 
codium               2594   6.285 2.046 /usr/share/codium/codium
kwin_x11             1821   4.487 1.231 /usr/bin/kwin_x11 --replace
brave                3370   4.459 1.630 /opt/brave.com/brave/brave
...

```

| Campo | Descripción |
|-------|-------------|
| Local Address | Dirección IP y puerto (formato IP:Puerto) |
| Proto | Protocolo de red (TCP, UDP, UNIX) |
| PID | Identificador del proceso |
| Process | Nombre del proceso que mantiene el puerto |

---

## Arquitectura del Proyecto

### Estructura de Directorios

```
peek/
├── main.go          # Punto de entrada, importa cli.Exec()
├── go.mod           # Definición de dependencias Go
├── go.sum           # Checksums de dependencias
│
└── cli/             # Paquete principal de CLI
    ├── root.go      # Definición del comando raíz "peek"
    ├── commands.go  # Registro de subcomandos (cpu, mem, stg, port)
    │
    └── cmds/        # Paquete de implementación de comandos
        ├── cpu.go   # Lógica del comando 'cpu'
        ├── mem.go   # Lógica del comando 'mem'
        ├── stg.go   # Lógica del comando 'stg'
        ├── port.go  # Lógica del comando 'port'
        ├── proc.go  # Lógica del comando 'proc'
        └── ops.go   # Utilidades compartidas (humanSize)
```

### Flujo de Ejecución

```
┌─────────────────────────────────────────────────────────────┐
│                        main.go                              │
│  Llama a cli.Exec() para iniciar la aplicación              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                         cli.Exec()                          │
│  Ejecuta peekCmd.Run() con contexto y argumentos            │
│  Maneja errores fatales con log.Fatal()                     │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  peekCmd (root.go)                          │
│  Nombre: "peek"                                             │
│  Usage: "Just a simple peek app..."                         │
│  Commands: [cpuPeek, memPeek, stgPeek, portPeek]            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              Command Handler (commands.go)                  │
│  cpuPeek   → cmds.CPU()                                     │
│  memPeek   → cmds.MEM()                                     │
│  stgPeek   → cmds.STG()                                     │
│  portPeek  → cmds.PORT(portFlag)                            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Implementation (cmds/)                     │
│                                                             │
│  cpu.go  → gopsutil/cpu + gopsutil/host                     │
│  mem.go  → gopsutil/mem                                     │
│  stg.go  → gopsutil/disk                                    │
│  port.go → gopsutil/net + gopsutil/process                  │
│  ops.go  → Formato de bytes (humanSize)                     │
└─────────────────────────────────────────────────────────────┘
```

---

## Dependencias

Peek utiliza dos dependencias principales:

### 1. urfave/cli/v3

**Repositorio:** https://github.com/urfave/cli  
**Documentación:** https://cli.urfave.org

Framework declarativo para construir aplicaciones CLI en Go.

**Características utilizadas:**
- Registro de comandos y subcomandos
- Manejo de argumentos posicionales
- Sistema de ayuda automático
- Completado de shell (bash, zsh, fish, powershell)

### 2. shirou/gopsutil/v4

**Repositorio:** https://github.com/shirou/gopsutil  
**Documentación:** https://pkg.go.dev/github.com/shirou/gopsutil/v4

Puerto en Go de la librería Python `psutil`. Proporciona acceso a información del sistema.

**Módulos utilizados:**

| Módulo | Paquete | Uso |
|--------|---------|-----|
| `cpu` | `github.com/shirou/gopsutil/v4/cpu` | Información del CPU (modelo, frecuencia, cores) |
| `host` | `github.com/shirou/gopsutil/v4/host` | Información del SO (OS, kernel, plataforma) |
| `mem` | `github.com/shirou/gopsutil/v4/mem` | Estadísticas de memoria virtual |
| `disk` | `github.com/shirou/gopsutil/v4/disk` | Información de particiones y uso de disco |
| `net` | `github.com/shirou/gopsutil/v4/net` | Conexiones de red y sockets |
| `process` | `github.com/shirou/gopsutil/v4/process` | Información de procesos por PID |

---

## Guía de Desarrollo

### Añadir un Nuevo Comando

Supongamos que quieres añadir un comando `net` para ver estadísticas de red:

**Paso 1: Crear el archivo de implementación**

```go
// cmds/net.go
package cmds

import (
    "fmt"
    "log"
    "github.com/shirou/gopsutil/v4/net"
)

func NET() {
    counters, err := net.IOCounters()
    if err != nil {
        log.Fatalf("could not get network stats: %v", err)
    }

    for _, c := range counters {
        ....
```

**Paso 2: Registrar el comando en `commands.go`**

```go
var netPeek = &cli.Command{
    Name:  "net",
    Usage: "peek net",
    Action: func(ctx context.Context, cmd *cli.Command) error {
        cmds.NET()
        return nil
    },
}
```

**Paso 3: Añadir al comando raíz en `root.go`**

```go
var peekCmd = &cli.Command{
    Name:  "peek",
    Usage: "Just a simple peek app that gives me information about linux",
    Commands: []*cli.Command{
        cpuPeek,
        memPeek,
        stgPeek,
        portPeek,
        netPeek,  // ← Añadir aquí
    },
}
```

### Convenciones de Código

| Convención | Ejemplo | Descripción |
|------------|---------|-------------|
| Nombres de funciones | `PORT()`, `CPU()`, `MEM()` | Mayúsculas, alineación con comandos CLI |
| Estructuras de datos | `connectionInfo` | CamelCase sin acrónimos |
| Constantes | `KB`, `MB`, `GB` | Mayúsculas completas |
| Manejo de errores | `log.Fatalf()` en init, retorna error en action | Fail-fast en setup, propagación en action |
| Salida | `fmt.Printf` con formatos alineados | Tablas con anchos fijos para legibilidad |

---

## Ejemplos de Salida

### Ejemplo 1: CPU

```
OS: Linux - Platform ubuntu
Kernel: 5.15.0 - x86_64
CPU: 11th Gen Intel(R) Core(TM) i7-11700 @ 2.50GHz - Cores 16
```

### Ejemplo 2: Memoria

```
Total: 31.17 GiB, Free: 18.52 GiB, Used: 12.65 GiB, UsedPercent: 40.584%
```

### Ejemplo 3: Almacenamiento

```
MOUNT                 TYPE      TOTAL     USED     FREE
/                     ext4      49.97 GB  24.50 GB 25.47 GB
/boot/efi             vfat      1.00 GB   0.01 GB  0.99 GB
/home                 ext4      98.51 GB  45.00 GB 53.51 GB

TOTAL STORAGE OVERVIEW
Total: 149.48 GB  Used: 69.51 GB  Free: 79.97 GB
```

### Ejemplo 4: Puertos

```
Local Address        Proto  PID    Process
0.0.0.0:22           TCP    892    sshd
127.0.0.1:631        TCP    1234   cupsd
0.0.0.0:8080         TCP    5678   java
:::33060             TCP    9012   mysqld
```

---

## Licencia

Este proyecto es software libre bajo la licencia GNU GENERAL PUBLIC LICENSE Version 3. Consulta el archivo `LICENSE` para más detalles.

---

*Generado automáticamente para el proyecto peek v1.0.0*