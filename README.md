# Peek — Información del sistema desde la terminal

## Tabla de Contenidos

1. [Descripción](#descripción)
2. [Requisitos](#requisitos)
3. [Instalación](#instalación)
4. [Uso](#uso)
   - [peek cpu](#peek-cpu)
   - [peek mem](#peek-mem)
   - [peek stg](#peek-stg)
   - [peek port](#peek-port)
   - [peek proc](#peek-proc)
   - [Salida JSON](#salida-json)
5. [Arquitectura del Proyecto](#arquitectura-del-proyecto)
   - [Estructura de Directorios](#estructura-de-directorios)
   - [Flujo de Ejecución](#flujo-de-ejecución)
6. [Dependencias](#dependencias)
7. [Guía de Desarrollo](#guía-de-desarrollo)
   - [Añadir un Nuevo Comando](#añadir-un-nuevo-comando)
   - [Convenciones de Código](#convenciones-de-código)
8. [Ideas y funcionalidades Unix solicitadas con frecuencia](#ideas-y-funcionalidades-unix-solicitadas-con-frecuencia)
9. [Limitaciones](#limitaciones)
10. [Licencia](#licencia)

---

## Descripción

**Peek** es una herramienta de línea de comandos escrita en **Go** para consultar información del sistema y sus procesos. Usa `gopsutil` para leer CPU, memoria, almacenamiento, conexiones y procesos. Aunque nació como una herramienta para Linux, sus dependencias admiten varios sistemas operativos; el soporte real depende de las funciones disponibles en cada plataforma.

```
peek cpu             # Información del procesador y del sistema operativo
peek mem             # Estadísticas de memoria virtual
peek stg             # Uso de almacenamiento por punto de montaje
peek port            # Sockets TCP en escucha
peek port 8080       # Filtrar sockets por puerto
peek proc            # Procesos y métricas de CPU/memoria
peek proc cpu        # 10 procesos con mayor uso de CPU
peek proc mem        # 10 procesos con mayor uso de memoria
peek proc pid 1234   # Detalles de un proceso
peek --help          # Ayuda de la CLI
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
git clone https://github.com/navigator-systems/peek.git
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

# Compilación optimizada
go build -ldflags="-s -w" -o peek .
```

### 4. Ejecutar

```bash
# Ejecutar directamente con Go
go run . --help

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
| Kernel | Versión del kernel |
| KernelArch | Arquitectura del kernel (x86_64, arm64, etc.) |
| CPU | Modelo completo del procesador |
| Cores | Cantidad de núcleos lógicos disponibles |

El formato JSON también está disponible con `--json`; consulta [Salida JSON](#salida-json).

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

El formato JSON también está disponible con `--json`.

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

El comando `stg` no tiene salida JSON todavía.

---

### `peek port`

Lista sockets TCP en estado `LISTEN` (esperando conexiones). El acceso al PID y al nombre del proceso puede depender de los permisos del usuario.

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
| Proto | Protocolo del socket (en esta lista, TCP) |
| PID | Identificador del proceso |
| Process | Nombre del proceso que mantiene el puerto |

Usa `peek --json port` para obtener JSON.

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

# Ver detalles de un proceso por PID
peek proc pid 1234
```

Las métricas de CPU se calculan a partir de dos lecturas separadas por un segundo. En sistemas multinúcleo, un proceso puede superar el 100% si consume más de un núcleo. Los comandos `cpu` y `mem` muestran como máximo diez procesos, en orden descendente. Para obtener JSON, agrega `--json`, por ejemplo `peek --json proc cpu`.

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
| Name | Nombre del proceso |
| PID | Identificador del proceso |
| CPU% | Porcentaje de CPU observado durante el intervalo de muestreo |
| Mem% | Porcentaje de memoria física usado por el proceso |
| Command | Línea de comandos del proceso |

El modo `pid` muestra nombre, PID, usuario, estado y comando. Un PID inválido o que ya no existe se informa como error.

### Salida JSON

`--json` es una opción global y se admite en `cpu`, `mem`, `port` y `proc`. Ejemplos:

```bash
peek --json cpu
peek --json mem
peek --json port 8080
peek --json proc
peek --json proc pid 1234
```

La salida JSON está pensada para scripts y herramientas como `jq`. `stg` actualmente solo produce una tabla de texto.

---

## Arquitectura del Proyecto

### Estructura de Directorios

```
peek/
├── main.go          # Punto de entrada
├── go.mod           # Módulo y dependencias
├── cli/             # Definición y registro de comandos
│   ├── root.go
│   ├── commands.go
│   └── flags.go
└── cmds/            # Implementación de comandos y utilidades
    ├── cpu.go
    ├── mem.go
    ├── stg.go
    ├── port.go
    ├── proc.go
    └── ops.go
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
│  Comandos: cpu, mem, stg, port, proc                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              Command Handler (commands.go)                  │
│  cpuPeek   → cmds.CPU()                                     │
│  memPeek   → cmds.MEM()                                     │
│  stgPeek   → cmds.STG()                                     │
│  portPeek  → cmds.PORT()                                    │
│  procPeek  → cmds.PROC()                                    │
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
│  proc.go → gopsutil/process                                 │
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
| `process` | `github.com/shirou/gopsutil/v4/process` | Información y métricas de procesos |

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
| Manejo de errores | Retorna errores desde la acción del comando | La CLI los informa y evita terminar el proceso desde la lógica de negocio |
| Salida | `fmt.Printf` con formatos alineados | Tablas con anchos fijos para legibilidad |

---

## Ideas y funcionalidades Unix solicitadas con frecuencia

Estas son ideas para ampliar Peek; todavía no forman parte de los comandos disponibles:

| Idea | Equivalente Unix frecuente | Valor para Peek |
|------|----------------------------|-----------------|
| Interfaces, direcciones IP, rutas y contadores de tráfico | `ip`, `ifconfig`, `ip -s link` | Diagnosticar conectividad desde `peek net` |
| Tiempo activo, carga y usuarios conectados | `uptime`, `who` | Resumen rápido del estado general, por ejemplo `peek uptime` |
| Procesos con filtros por nombre, usuario o PID | `ps`, `pgrep`, `top` | Buscar y ordenar procesos sin canalizaciones externas |
| Uso de disco por directorio | `du`, además de `df` | Encontrar directorios grandes; requiere una ruta y límites claros |
| Sockets establecidos y UDP, no solo TCP en escucha | `ss`, `netstat`, `lsof` | Filtrar por protocolo, estado, dirección y proceso |
| Servicios y logs recientes | `systemctl`, `journalctl` | Atajos para revisar el estado y mensajes de un servicio |
| Actualización periódica | `watch` | Opción `--watch` con intervalo configurable y salida legible en terminal |
| Salida JSON consistente para todos los comandos | `jq` en scripts | JSON también para almacenamiento y un esquema documentado |
| Señales de proceso (terminar, pausar, continuar) | `kill`, `pkill` | Solo con confirmación, validación de PID y advertencias; Peek actualmente no modifica procesos |

Como optimizaciones adicionales, conviene evitar consultas repetidas a la misma partición, devolver errores desde los comandos en vez de terminar el proceso con `log.Fatal`, agregar pruebas de formatos y ordenamiento, y comprobar cada función en Linux, macOS y FreeBSD antes de prometer compatibilidad.

## Limitaciones

- El listado de procesos muestrea CPU durante un segundo; procesos que terminan o dejan de ser accesibles durante el muestreo pueden no aparecer.
- Los permisos del sistema pueden limitar la información de procesos y sockets.
- `stg` suma el uso de todos los puntos de montaje visibles; montajes enlazados o superpuestos pueden hacer que el total no represente capacidad física única.
- `stg` aún no admite JSON; `port` lista únicamente sockets TCP en escucha.
- El soporte concreto de cada función depende del sistema operativo, su configuración y los permisos del usuario.

---

## Licencia

Este proyecto es software libre bajo la licencia GNU GENERAL PUBLIC LICENSE Version 3. Consulta el archivo `LICENSE` para más detalles.

---

*Generado automáticamente para el proyecto peek v1.0.0*
