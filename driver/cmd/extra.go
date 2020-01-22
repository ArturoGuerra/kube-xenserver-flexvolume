func initCommand() {
    capabilities := driverCapabilities{
        Attach: true,
    }

    output := driverOutput{
        Status: "Success",
        Capabilities: capabilities,
    }

    printLn(output)
}

func notSupported() {
    output := driverOutput{
        Status: "Not supported",
    }

    printLn(output)
}
