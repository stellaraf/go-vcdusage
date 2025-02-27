![stellar](https://res.cloudinary.com/stellaraf/image/upload/v1604277355/stellar-logo-gradient.png?width=300)

## `go-vcdusage`

[![Go Reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://pkg.go.dev/go.stellar.af/go-vcdusage) [![GitHub Tag](https://img.shields.io/github/v/tag/stellaraf/go-vcdusage?style=for-the-badge&label=Version)](https://github.com/stellaraf/go-vcdusage/tags)

## Testing Environment Variables

The following environment variables must be set (and valid) for tests to run.

| Key            | Description                              |
| :------------- | :--------------------------------------- |
| `VCD_URL`      | vCloud URL                               |
| `VCD_USERNAME` | vCloud Username                          |
| `VCD_PASSWORD` | vCloud Password                          |
| `VCD_ORG_ID`   | vCloud Org ID to test against            |
| `VCD_CORES`    | Number of allocated cores in vCloud Org  |
| `VCD_MEMORY`   | Amount of allocated memory in vCloud Org |
| `VCD_STORAGE`  | Amount allocated storage in vCloud Org   |
| `VCD_VM_COUNT` | Number of VMs in vCloud Org              |

