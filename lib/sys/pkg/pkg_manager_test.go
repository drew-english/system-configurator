package pkg_test

import (
	"fmt"
	"testing"

	"github.com/drew-english/system-configurator/internal/model"
	"github.com/drew-english/system-configurator/lib/sys"
	"github.com/drew-english/system-configurator/lib/sys/pkg"
	"github.com/drew-english/system-configurator/spec/stub/run"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PackageManager", func() {
	var (
		commandStubs     *run.CommandStubManager
		teardownCmdStubs func(testing.TB)
		manager          pkg.PacakgeManager
	)

	BeforeEach(func() {
		commandStubs, teardownCmdStubs = run.StubCommand()

		unregister := run.StubFind(sys.SupportedPackageManagers()[0], nil)
		defer unregister()

		var err error
		manager, err = pkg.FindPackageManager()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		teardownCmdStubs(GinkgoTB())
	})

	Describe("Name", func() {
		It("returns the base command", func() {
			Expect(manager.Name()).To(Equal(sys.SupportedPackageManagers()[0]))
		})
	})

	Describe("AddPackage", func() {
		pkg := &model.Package{Name: "some-test-pkg", Version: "1.2.3"}
		var addPkgExpression string

		BeforeEach(func() {
			addPkgExpression = fmt.Sprintf("%s (add|install|-S).* %s.*%s", manager.Name(), pkg.Name, pkg.Version)
		})

		It("returns nil", func() {
			commandStubs.Register(addPkgExpression, "package added successfully")
			Expect(manager.AddPackage(pkg)).ToNot(HaveOccurred())
		})

		Context("when the command fails", func() {
			It("returns an error", func() {
				commandStubs.RegisterError(addPkgExpression, 1, "failed to find package")
				Expect(manager.AddPackage(pkg)).To(MatchError(fmt.Sprintf("failed to find package\n%s: %s", manager.Name(), "generic error")))
			})
		})
	})

	Describe("RemovePackage", func() {
		pkg := &model.Package{Name: "some-test-pkg", Version: "1.2.3"}
		var rmPkgExpression string

		BeforeEach(func() {
			rmPkgExpression = fmt.Sprintf("%s (del|remove|erase|uninstall|-Rscn).* %s", manager.Name(), pkg.Name)
		})

		It("returns nil", func() {
			commandStubs.Register(rmPkgExpression, "package removed successfully")
			Expect(manager.RemovePackage(pkg)).ToNot(HaveOccurred())
		})

		Context("when the command fails", func() {
			It("returns an error", func() {
				commandStubs.RegisterError(rmPkgExpression, 1, "failed to find package")
				Expect(manager.RemovePackage(pkg)).To(MatchError(fmt.Sprintf("failed to find package\n%s: %s", manager.Name(), "generic error")))
			})
		})
	})

	Describe("ListPackages", func() {
		stubbedListOutput := map[string]string{
			"apk": `WARNING: opening from cache https://dl-cdn.alpinelinux.org/alpine/v3.19/main: No such file or directory
WARNING: opening from cache https://dl-cdn.alpinelinux.org/alpine/v3.19/community: No such file or directory
alpine-baselayout-3.4.3-r2 aarch64 {alpine-baselayout} (GPL-2.0-only) [installed]
alpine-baselayout-data-3.4.3-r2 aarch64 {alpine-baselayout} (GPL-2.0-only) [installed]
alpine-keys-2.4-r1 aarch64 {alpine-keys} (MIT) [installed]
apk-tools-2.14.0-r5 aarch64 {apk-tools} (GPL-2.0-only) [installed]
busybox-1.36.1-r15 aarch64 {busybox} (GPL-2.0-only) [installed]
busybox-binsh-1.36.1-r15 aarch64 {busybox} (GPL-2.0-only) [installed]
ca-certificates-bundle-20230506-r0 aarch64 {ca-certificates} (MPL-2.0 AND MIT) [installed]
libc-utils-0.7.2-r5 aarch64 {libc-dev} (BSD-2-Clause AND BSD-3-Clause) [installed]
libcrypto3-3.1.4-r5 aarch64 {openssl} (Apache-2.0) [installed
ssl_client-1.36.1-r15 aarch64 {busybox} (GPL-2.0-only) [installed]`,
			"apt": `Listing... Done
adduser/now 3.118ubuntu5 all [installed,local]
apt/now 2.4.11 arm64 [installed,local]
base-files/now 12ubuntu4.5 arm64 [installed,local]
base-passwd/now 3.5.52build1 arm64 [installed,local]
bash/now 5.1-6ubuntu1 arm64 [installed,local]
bsdutils/now 1:2.37.2-4ubuntu3 arm64 [installed,local]
coreutils/now 8.32-4.1ubuntu1.1 arm64 [installed,local]
dash/now 0.5.11+git20210903+057cd650a4ed-3build1 arm64 [installed,local]
debconf/now 1.5.79ubuntu1 all [installed,local]
debianutils/now 5.5-1ubuntu2 arm64 [installed,local]
diffutils/now 1:3.8-0ubuntu2 arm64 [installed,local]
dpkg/now 1.21.1ubuntu2.2 arm64 [installed,local]`,
			"brew": `abseil 20230802.1
aom 3.8.1
argocd 2.10.1
aribb24 1.0.4
autoconf 2.72
automake 1.16.5
aws-sam-cli 1.109.0_1
bdw-gc 8.2.6
boost 1.84.0
brotli 1.1.0
c-ares 1.26.0
ca-certificates 2023-12-12
cairo 1.18.0
cjson 1.7.17`,
			"dnf": `alternatives.aarch64                                                      1.26-1.fc39                                                 @koji-override-1
audit-libs.aarch64                                                        3.1.2-8.fc39                                                @koji-override-1
authselect.aarch64                                                        1.4.3-1.fc39                                                @anaconda
authselect-libs.aarch64                                                   1.4.3-1.fc39                                                @anaconda
basesystem.noarch                                                         11-18.fc39                                                  @anaconda
bash.aarch64                                                              5.2.26-1.fc39                                               @koji-override-1
bzip2-libs.aarch64                                                        1.0.8-16.fc39                                               @anaconda
ca-certificates.noarch                                                    2023.2.60_v7.0.306-2.fc39                                   @anaconda
coreutils.aarch64                                                         9.3-5.fc39                                                  @koji-override-1
coreutils-common.aarch64                                                  9.3-5.fc39                                                  @koji-override-1`,
			// "snap":   ``,
			"pacman": `warning: database file for 'core' does not exist (use '-Sy' to download)
warning: database file for 'extra' does not exist (use '-Sy' to download)
warning: database file for 'alarm' does not exist (use '-Sy' to download)
warning: database file for 'aur' does not exist (use '-Sy' to download)
acl 2.3.2-1
archlinux-keyring 20240208-1
archlinuxarm-keyring 20140119-2
argon2 20190702-5
attr 2.5.2-1
audit 4.0-1
base 3-2
bash 5.2.026-2
brotli 1.1.0-1
bzip2 1.0.8-5
ca-certificates 20220905-1
ca-certificates-mozilla 3.98-1
ca-certificates-utils 20220905-1
coreutils 9.4-3
cryptsetup 2.7.0-1
curl 8.6.0-3
dbus 1.14.10-2
dbus-broker 35-2
dbus-broker-units 35-2
device-mapper 2.03.23-1
e2fsprogs 1.47.0-1`,
		}

		expectedPkgList := map[string][]*model.Package{
			"apk":  {{Name: "busybox", Version: "1.36.1-r15"}, {Name: "ssl_client", Version: "1.36.1-r15"}, {Name: "ca-certificates-bundle", Version: "20230506-r0"}},
			"apt":  {{Name: "adduser", Version: "3.118ubuntu5"}, {Name: "apt", Version: "2.4.11"}, {Name: "diffutils", Version: "1:3.8-0ubuntu2"}},
			"brew": {{Name: "argocd", Version: "2.10.1"}, {Name: "ca-certificates", Version: "2023-12-12"}, {Name: "cairo", Version: "1.18.0"}},
			"dnf":  {{Name: "alternatives", Version: "1.26"}, {Name: "authselect", Version: "1.4.3"}, {Name: "basesystem", Version: "11"}},
			// "snap":   {},
			"pacman": {{Name: "acl", Version: "2.3.2-1"}, {Name: "archlinux-keyring", Version: "20240208-1"}, {Name: "argon2", Version: "20190702-5"}},
		}

		for mgrName, mgr := range pkg.Managers {
			Context(fmt.Sprintf("when the package manager is %s", mgrName), func() {
				BeforeEach(func() {
					manager = mgr
					commandStubs.Register(fmt.Sprintf("%s (list|-Q)( --installed){0,1}", mgrName), stubbedListOutput[mgrName])
				})

				It("should return the list of packages", func() {
					pkgList, err := manager.ListPackages()
					Expect(err).NotTo(HaveOccurred())
					Expect(pkgList).To(ContainElements(expectedPkgList[mgrName]))
				})
			})
		}
	})
})
