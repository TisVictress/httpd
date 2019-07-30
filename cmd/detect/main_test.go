/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/httpd-cnb/httpd"
	"github.com/cloudfoundry/php-dist-cnb/php"

	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitDetect(t *testing.T) {
	spec.Run(t, "Detect", testDetect, spec.Report(report.Terminal{}))
}

func testDetect(t *testing.T, when spec.G, it spec.S) {
	var factory *test.DetectFactory

	it.Before(func() {
		RegisterTestingT(t)
		factory = test.NewDetectFactory(t)
	})

	when("there is an httpd.conf", func() {
		it.Before(func() {
			test.WriteFile(t, filepath.Join(factory.Detect.Application.Root, "httpd.conf"), "")
		})

		it("should pass with the default version of httpd", func() {
			code, err := runDetect(factory.Detect)
			Expect(err).NotTo(HaveOccurred())

			Expect(code).To(Equal(detect.PassStatusCode))

			Expect(factory.Output).To(Equal(buildplan.BuildPlan{
				httpd.Dependency: buildplan.Dependency{
					Metadata: buildplan.Metadata{"launch": true},
				},
			}))
		})

		when("there is a buildpack.yml", func() {
			it("should request the supplied version", func() {
				yaml := "{'httpd': {'version': 1.2.3}}"
				test.WriteFile(t, filepath.Join(factory.Detect.Application.Root, "buildpack.yml"), yaml)

				code, err := runDetect(factory.Detect)
				Expect(err).NotTo(HaveOccurred())

				Expect(code).To(Equal(detect.PassStatusCode))

				Expect(factory.Output).To(Equal(buildplan.BuildPlan{
					httpd.Dependency: buildplan.Dependency{
						Version:  "1.2.3",
						Metadata: buildplan.Metadata{"launch": true},
					},
				}))
			})

			it("should request the default version when no version is requested", func() {
				test.WriteFile(t, filepath.Join(factory.Detect.Application.Root, "buildpack.yml"), " ")

				code, err := runDetect(factory.Detect)
				Expect(err).NotTo(HaveOccurred())

				Expect(code).To(Equal(detect.PassStatusCode))

				Expect(factory.Output).To(Equal(buildplan.BuildPlan{
					httpd.Dependency: buildplan.Dependency{
						Metadata: buildplan.Metadata{"launch": true},
					},
				}))
			})
		})
	})

	when("there is NOT an httpd.conf", func() {
		it("should pass if `php-binary` is in the build plan", func() {
			factory.AddBuildPlan(php.Dependency, buildplan.Dependency{})

			code, err := runDetect(factory.Detect)
			Expect(err).NotTo(HaveOccurred())

			Expect(code).To(Equal(detect.PassStatusCode))
			Expect(factory.Output).To(Equal(buildplan.BuildPlan{}))
		})

		it("should fail if `php-binary` is not in the build plan", func() {
			code, err := runDetect(factory.Detect)
			Expect(err).To(HaveOccurred())

			Expect(code).To(Equal(detect.FailStatusCode))
		})
	})
}
