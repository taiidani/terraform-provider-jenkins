package jenkins

import (
	"encoding/xml"
	"reflect"
	"strings"
	"testing"
)

func Test_parseFolder(t *testing.T) {
	type args struct {
		def string
	}
	tests := []struct {
		name    string
		args    args
		want    *folder
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				def: `<com.cloudbees.hudson.plugins.folder.Folder plugin="cloudbees-folder@6.15">
  <actions/>
  <description>Example Description</description>
  <displayName>Example Display Name</displayName>
  <properties>
    <com.cloudbees.hudson.plugins.folder.properties.AuthorizationMatrixProperty>
      <inheritanceStrategy class="org.jenkinsci.plugins.matrixauth.inheritance.InheritParentStrategy"/>
      <permission>com.cloudbees.plugins.credentials.CredentialsProvider.Create:anonymous</permission>
      <permission>com.cloudbees.plugins.credentials.CredentialsProvider.Delete:authenticated</permission>
      <permission>hudson.model.Item.Cancel:authenticated</permission>
      <permission>hudson.model.Item.Discover:anonymous</permission>
    </com.cloudbees.hudson.plugins.folder.properties.AuthorizationMatrixProperty>
    <org.jenkinsci.plugins.workflow.libs.FolderLibraries plugin="workflow-cps-global-lib@2.17">
      <libraries>
        <org.jenkinsci.plugins.workflow.libs.LibraryConfiguration>
          <name>Example Library Configuration</name>
          <implicit>false</implicit>
          <allowVersionOverride>true</allowVersionOverride>
          <includeInChangesets>true</includeInChangesets>
        </org.jenkinsci.plugins.workflow.libs.LibraryConfiguration>
      </libraries>
    </org.jenkinsci.plugins.workflow.libs.FolderLibraries>
  </properties>
  <folderViews class="com.cloudbees.hudson.plugins.folder.views.DefaultFolderViewHolder">
    <views>
      <hudson.model.AllView>
        <owner class="com.cloudbees.hudson.plugins.folder.Folder" reference="../../../.."/>
        <name>All</name>
        <filterExecutors>false</filterExecutors>
        <filterQueue>false</filterQueue>
        <properties class="hudson.model.View$PropertyList"/>
      </hudson.model.AllView>
      <hudson.model.ListView>
        <owner class="com.cloudbees.hudson.plugins.folder.Folder" reference="../../../.."/>
        <name>Example View</name>
        <filterExecutors>false</filterExecutors>
        <filterQueue>false</filterQueue>
        <properties class="hudson.model.View$PropertyList"/>
        <jobNames>
          <comparator class="hudson.util.CaseInsensitiveComparator"/>
        </jobNames>
        <jobFilters/>
        <columns>
          <hudson.views.StatusColumn/>
          <hudson.views.WeatherColumn/>
          <hudson.views.JobColumn/>
          <hudson.views.LastSuccessColumn/>
          <hudson.views.LastFailureColumn/>
          <hudson.views.LastDurationColumn/>
          <hudson.views.BuildButtonColumn/>
        </columns>
        <recurse>false</recurse>
      </hudson.model.ListView>
    </views>
    <primaryView>All</primaryView>
    <tabBar class="hudson.views.DefaultViewsTabBar"/>
  </folderViews>
  <healthMetrics>
    <com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric>
      <nonRecursive>true</nonRecursive>
    </com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric>
  </healthMetrics>
  <icon class="com.cloudbees.hudson.plugins.folder.icons.StockFolderIcon"/>
</com.cloudbees.hudson.plugins.folder.Folder>`,
			},
			want: &folder{
				XMLName:     xml.Name{Local: "com.cloudbees.hudson.plugins.folder.Folder"},
				Description: "Example Description",
				DisplayName: "Example Display Name",
				Properties: folderProperties{
					Security: &folderSecurity{
						InheritanceStrategy: folderPermissionInheritanceStrategy{
							Class: "org.jenkinsci.plugins.matrixauth.inheritance.InheritParentStrategy",
						},
						Permission: []string{
							"com.cloudbees.plugins.credentials.CredentialsProvider.Create:anonymous",
							"com.cloudbees.plugins.credentials.CredentialsProvider.Delete:authenticated",
							"hudson.model.Item.Cancel:authenticated",
							"hudson.model.Item.Discover:anonymous",
						},
					},
					Other: []xmlRawProperty{
						{
							XMLName: xml.Name{Local: "org.jenkinsci.plugins.workflow.libs.FolderLibraries"},
							Plugin:  "workflow-cps-global-lib@2.17",
							Raw: `
      <libraries>
        <org.jenkinsci.plugins.workflow.libs.LibraryConfiguration>
          <name>Example Library Configuration</name>
          <implicit>false</implicit>
          <allowVersionOverride>true</allowVersionOverride>
          <includeInChangesets>true</includeInChangesets>
        </org.jenkinsci.plugins.workflow.libs.LibraryConfiguration>
      </libraries>
    `,
						},
					},
				},
				FolderViews: xmlRawProperty{
					XMLName: xml.Name{Local: "folderViews"},
					Raw: `
    <views>
      <hudson.model.AllView>
        <owner class="com.cloudbees.hudson.plugins.folder.Folder" reference="../../../.."/>
        <name>All</name>
        <filterExecutors>false</filterExecutors>
        <filterQueue>false</filterQueue>
        <properties class="hudson.model.View$PropertyList"/>
      </hudson.model.AllView>
      <hudson.model.ListView>
        <owner class="com.cloudbees.hudson.plugins.folder.Folder" reference="../../../.."/>
        <name>Example View</name>
        <filterExecutors>false</filterExecutors>
        <filterQueue>false</filterQueue>
        <properties class="hudson.model.View$PropertyList"/>
        <jobNames>
          <comparator class="hudson.util.CaseInsensitiveComparator"/>
        </jobNames>
        <jobFilters/>
        <columns>
          <hudson.views.StatusColumn/>
          <hudson.views.WeatherColumn/>
          <hudson.views.JobColumn/>
          <hudson.views.LastSuccessColumn/>
          <hudson.views.LastFailureColumn/>
          <hudson.views.LastDurationColumn/>
          <hudson.views.BuildButtonColumn/>
        </columns>
        <recurse>false</recurse>
      </hudson.model.ListView>
    </views>
    <primaryView>All</primaryView>
    <tabBar class="hudson.views.DefaultViewsTabBar"/>
  `,
				},
				HealthMetrics: xmlRawProperty{
					XMLName: xml.Name{Local: "healthMetrics"},
					Raw: `
    <com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric>
      <nonRecursive>true</nonRecursive>
    </com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric>
  `,
				},
			},
		},
		{
			name: "error-invalid-xml",
			args: args{
				def: `Invalid`,
			},
			want:    &folder{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFolder(tt.args.def)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFolder() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_folder_Render(t *testing.T) {
	type fields struct {
		Description string
		DisplayName string
		Properties  folderProperties
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Description: "Example Description",
				DisplayName: "Example Display Name",
				Properties: folderProperties{
					Security: &folderSecurity{
						Permission: []string{"example", "permission"},
						InheritanceStrategy: folderPermissionInheritanceStrategy{
							Class: "org.jenkinsci.plugins.matrixauth.inheritance.InheritParentStrategy",
						},
					},
					Other: []xmlRawProperty{
						{
							XMLName: xml.Name{Local: "org.jenkinsci.plugins.workflow.libs.FolderLibraries"},
							Plugin:  "workflow-cps-global-lib@2.17",
							Raw: `
      <libraries>
        <org.jenkinsci.plugins.workflow.libs.LibraryConfiguration>
          <name>Example Library Configuration</name>
          <implicit>false</implicit>
          <allowVersionOverride>true</allowVersionOverride>
          <includeInChangesets>true</includeInChangesets>
        </org.jenkinsci.plugins.workflow.libs.LibraryConfiguration>
      </libraries>
    `,
						},
					},
				},
			},
			want: []byte(`<com.cloudbees.hudson.plugins.folder.Folder>
	<description>Example Description</description>
    <displayName>Example Display Name</displayName>
	<properties>
		<com.cloudbees.hudson.plugins.folder.properties.AuthorizationMatrixProperty>
      <inheritanceStrategy class="org.jenkinsci.plugins.matrixauth.inheritance.InheritParentStrategy"></inheritanceStrategy>
			<permission>example</permission>
			<permission>permission</permission>
    </com.cloudbees.hudson.plugins.folder.properties.AuthorizationMatrixProperty>
    <org.jenkinsci.plugins.workflow.libs.FolderLibraries plugin="workflow-cps-global-lib@2.17">
      <libraries>
        <org.jenkinsci.plugins.workflow.libs.LibraryConfiguration>
          <name>Example Library Configuration</name>
          <implicit>false</implicit>
          <allowVersionOverride>true</allowVersionOverride>
          <includeInChangesets>true</includeInChangesets>
        </org.jenkinsci.plugins.workflow.libs.LibraryConfiguration>
      </libraries>
    </org.jenkinsci.plugins.workflow.libs.FolderLibraries>
	</properties>
	<folderViews></folderViews>
	<healthMetrics></healthMetrics>
</com.cloudbees.hudson.plugins.folder.Folder>`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &folder{
				Description: tt.fields.Description,
				DisplayName: tt.fields.DisplayName,
				Properties:  tt.fields.Properties,
			}
			got, err := j.Render()
			if (err != nil) != tt.wantErr {
				t.Errorf("folder.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			strGot := strings.ReplaceAll(string(got), "\n", "")
			strGot = strings.ReplaceAll(strGot, "  ", "\t")
			strGot = strings.ReplaceAll(strGot, "\t", "")
			want := strings.ReplaceAll(string(tt.want), "\n", "")
			want = strings.ReplaceAll(want, "  ", "\t")
			want = strings.ReplaceAll(want, "\t", "")

			if strGot != want {
				t.Errorf("folder.Render() = %v, want %v", strGot, want)
			}
		})
	}
}
