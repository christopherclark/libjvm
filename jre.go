/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package libjvm

import (
	"fmt"
	"os"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
)

type JRE struct {
	Crush            crush.Crush
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
	Metadata         map[string]interface{}
}

func NewJRE(dependency libpak.BuildpackDependency, cache libpak.DependencyCache, metadata map[string]interface{},
	plan *libcnb.BuildpackPlan) JRE {

	return JRE{
		LayerContributor: libpak.NewDependencyLayerContributor(dependency, cache, plan),
		Logger:           bard.NewLogger(os.Stdout),
		Metadata:         metadata,
	}
}

func (j JRE) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Body("Expanding to %s", layer.Path)
		if err := j.Crush.ExtractTarGz(artifact, layer.Path, 1); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to expand JRE: %w", err)
		}

		layer.SharedEnvironment.Override("JAVA_HOME", layer.Path)
		layer.SharedEnvironment.Override("MALLOC_ARENA_MAX", "2")
		layer.Profile.Add("active-processor-count", `JAVA_OPTS="${JAVA_OPTS} -XX:ActiveProcessorCount=$(nproc)"
export JAVA_OPTS`)

		if v, ok := j.Metadata["build"].(bool); ok && v {
			layer.Build = true
			layer.Cache = true
		}

		if v, ok := j.Metadata["launch"].(bool); ok && v {
			layer.Launch = true
		}

		return layer, nil
	})
}

func (JRE) Name() string {
	return "jre"
}
