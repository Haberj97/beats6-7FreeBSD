// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// +build mage

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"

	auditbeat "github.com/elastic/beats/auditbeat/scripts/mage"
	"github.com/elastic/beats/dev-tools/mage"
)

func init() {
	mage.BeatDescription = "Audit the activities of users and processes on your system."
	mage.BeatLicense = "Elastic License"
	mage.Platforms = mage.Platforms.Filter("!linux/ppc64 !linux/mips64")
}

// Aliases provides compatibility with CI while we transition all Beats
// to having common testing targets.
var Aliases = map[string]interface{}{
	"goTestUnit": GoUnitTest, // dev-tools/jenkins_ci.ps1 uses this.
}

// Build builds the Beat binary.
func Build() error {
	return mage.Build(mage.DefaultBuildArgs())
}

// GolangCrossBuild build the Beat binary inside of the golang-builder.
// Do not use directly, use crossBuild instead.
func GolangCrossBuild() error {
	if d, ok := deps[mage.Platform.Name]; ok {
		mg.Deps(d)
	}
	return mage.GolangCrossBuild(mage.DefaultGolangCrossBuildArgs())
}

// CrossBuild cross-builds the beat for all target platforms.
func CrossBuild() error {
	return mage.CrossBuild()
}

// BuildGoDaemon builds the go-daemon binary (use crossBuildGoDaemon).
func BuildGoDaemon() error {
	return mage.BuildGoDaemon()
}

// CrossBuildGoDaemon cross-builds the go-daemon binary using Docker.
func CrossBuildGoDaemon() error {
	return mage.CrossBuildGoDaemon()
}

// Clean cleans all generated files and build artifacts.
func Clean() error {
	return mage.Clean()
}

// Package packages the Beat for distribution.
// Use SNAPSHOT=true to build snapshots.
// Use PLATFORMS to control the target platforms.
// Use VERSION_QUALIFIER to control the version qualifier.
func Package() {
	start := time.Now()
	defer func() { fmt.Println("package ran for", time.Since(start)) }()

	mage.UseElasticBeatXPackPackaging()
	mage.PackageKibanaDashboardsFromBuildDir()
	auditbeat.CustomizePackaging()

	mg.SerialDeps(Fields, Dashboards, Config, mage.GenerateModuleIncludeListGo)
	mg.Deps(CrossBuild, CrossBuildGoDaemon)
	mg.SerialDeps(mage.Package, TestPackages)
}

// TestPackages tests the generated packages (i.e. file modes, owners, groups).
func TestPackages() error {
	return mage.TestPackages()
}

// Fields generates a fields.yml and for each module generates a fields.go.
func Fields() {
	mg.SerialDeps(fieldsYML, moduleFieldsGo)
}

// moduleFieldsGo generates a fields.go file for each module.
func moduleFieldsGo() error {
	return mage.GenerateModuleFieldsGo("module")
}

// fieldsYML generates a fields.yml based on auditbeat + x-pack/auditbeat/modules.
func fieldsYML() error {
	return mage.GenerateFieldsYAML(mage.OSSBeatDir("module"), "module")
}

// ExportDashboard exports a dashboard and writes it into the correct directory.
//
// Required environment variables:
// - MODULE: Name of the module
// - ID:     Dashboard id
func ExportDashboard() error {
	return mage.ExportDashboard()
}

// Dashboards collects all the dashboards and generates index patterns.
func Dashboards() error {
	return mage.KibanaDashboards(mage.OSSBeatDir("module"), "module")
}

// Config generates both the short and reference configs.
func Config() error {
	return auditbeat.Config(
		mage.OSSBeatDir(auditbeat.ConfigTemplateGlob),
		auditbeat.ConfigTemplateGlob)
}

// Update is an alias for running fields, dashboards, config.
func Update() {
	mg.SerialDeps(Fields, Dashboards, Config, mage.GenerateModuleIncludeListGo,
		Docs)
}

// Docs collects the documentation.
func Docs() error {
	return auditbeat.CollectDocs(mage.OSSBeatDir(), auditbeat.XpackBeatDir())
}

// Fmt formats source code and adds file headers.
func Fmt() {
	mg.Deps(mage.Format)
}

// Check runs fmt and update then returns an error if any modifications are found.
func Check() {
	mg.SerialDeps(mage.Format, Update, mage.Check)
}

// IntegTest executes integration tests (it uses Docker to run the tests).
func IntegTest() {
	mage.AddIntegTestUsage()
	defer mage.StopIntegTestEnv()
	mg.SerialDeps(GoIntegTest, PythonIntegTest)
}

// UnitTest executes the unit tests.
func UnitTest() {
	mg.SerialDeps(GoUnitTest, PythonUnitTest)
}

// GoUnitTest executes the Go unit tests.
// Use TEST_COVERAGE=true to enable code coverage profiling.
// Use RACE_DETECTOR=true to enable the race detector.
func GoUnitTest(ctx context.Context) error {
	return mage.GoTest(ctx, mage.DefaultGoTestUnitArgs())
}

// GoIntegTest executes the Go integration tests.
// Use TEST_COVERAGE=true to enable code coverage profiling.
// Use RACE_DETECTOR=true to enable the race detector.
func GoIntegTest(ctx context.Context) error {
	return mage.RunIntegTest("goIntegTest", func() error {
		return mage.GoTest(ctx, mage.DefaultGoTestIntegrationArgs())
	})
}

// PythonUnitTest executes the python system tests.
func PythonUnitTest() error {
	mg.SerialDeps(Fields, mage.BuildSystemTestBinary)
	return mage.PythonNoseTest(mage.DefaultPythonTestUnitArgs())
}

// PythonIntegTest executes the python system tests in the integration environment (Docker).
func PythonIntegTest(ctx context.Context) error {
	if !mage.IsInIntegTestEnv() {
		mg.Deps(Fields)
	}
	return mage.RunIntegTest("pythonIntegTest", func() error {
		mg.Deps(mage.BuildSystemTestBinary)
		return mage.PythonNoseTest(mage.DefaultPythonTestIntegrationArgs())
	})
}

// -----------------------------------------------------------------------------
// - Install the librpm-dev package
var (
	deps = map[string]func() error{
		"linux/386":      installLinux386,
		"linux/amd64":    installLinuxAMD64,
		"linux/arm64":    installLinuxARM64,
		"linux/armv5":    installLinuxARMLE,
		"linux/armv6":    installLinuxARMLE,
		"linux/armv7":    installLinuxARMHF,
		"linux/mips":     installLinuxMIPS,
		"linux/mipsle":   installLinuxMIPSLE,
		"linux/mips64le": installLinuxMIPS64LE,
		"linux/ppc64le":  installLinuxPPC64LE,
		"linux/s390x":    installLinuxS390X,

		//"linux/ppc64":  installLinuxPpc64,
		//"linux/mips64": installLinuxMips64,
	}
)

const (
	librpmDevPkgName = "librpm-dev"
)

func installLinuxAMD64() error {
	return installDependencies(librpmDevPkgName, "")
}

func installLinuxARM64() error {
	return installDependencies(librpmDevPkgName+":arm64", "arm64")
}

func installLinuxARMHF() error {
	return installDependencies(librpmDevPkgName+":armhf", "armhf")
}

func installLinuxARMLE() error {
	return installDependencies(librpmDevPkgName+":armel", "armel")
}

func installLinux386() error {
	return installDependencies(librpmDevPkgName+":i386", "i386")
}

func installLinuxMIPS() error {
	return installDependencies(librpmDevPkgName+":mips", "mips")
}

func installLinuxMIPS64LE() error {
	return installDependencies(librpmDevPkgName+":mips64el", "mips64el")
}

func installLinuxMIPSLE() error {
	return installDependencies(librpmDevPkgName+":mipsel", "mipsel")
}

func installLinuxPPC64LE() error {
	return installDependencies(librpmDevPkgName+":ppc64el", "ppc64el")
}

func installLinuxS390X() error {
	return installDependencies(librpmDevPkgName+":s390x", "s390x")
}

func installDependencies(pkg, arch string) error {
	if arch != "" {
		err := sh.Run("dpkg", "--add-architecture", arch)
		if err != nil {
			return errors.Wrap(err, "error while adding architecture")
		}
	}

	// TODO: This is only for debian 7 and should be removed when move to a newer OS. This flag is
	// going to be used unnecessary when building using non-debian7 images
	// (like when making the linux/arm binaries) and we should remove it soonish.
	// See https://github.com/elastic/beats/issues/11750 for more details.
	if err := sh.Run("apt-get", "update", "-o", "Acquire::Check-Valid-Until=false"); err != nil {
		return err
	}

	return sh.Run("apt-get", "install", "-y", "--no-install-recommends", pkg)
}
