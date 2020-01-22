package main

import (
)

func umount(mountDir string) {
    var err error

    if err = detach(params); err != nil {
        failure(err)
    }

    if err = loadParamsFromFile(mountDir); err != nil {
        failure(err)
    }

    defer func() {
        if err = deleteParamsFile(mountDir); err != nil {
            failure(err)
        }
    }()

    if err = detach(params, mountDir); err != nil {
        failure(err)
    }

    success("Unmounted and detached volume")
}

func loadParamsFromFile(mountDir string) *jsonParams {
}

func deleteParamsFile(mountDir string) error {
}

func detach(params jsonParams) error {
}
