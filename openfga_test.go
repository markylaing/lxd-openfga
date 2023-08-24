package openfga

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/openfga/go-sdk/client"
	"github.com/stretchr/testify/suite"
)

type openFGASuite struct {
	suite.Suite
	fga         *client.OpenFgaClient
	authModelID *string
}

func TestOpenFGASuite(t *testing.T) {
	suite.Run(t, new(openFGASuite))
}

func (s *openFGASuite) SetupTest() {
	var err error
	s.fga, err = client.NewSdkClient(&client.ClientConfiguration{
		ApiScheme: "http",
		ApiHost:   "localhost:8080",
	})
	s.Require().NoError(err)

	createStoreResponse, err := s.fga.CreateStore(context.Background()).Body(client.ClientCreateStoreRequest{Name: "test"}).Execute()
	s.Require().NoError(err)

	s.fga.SetStoreId(*createStoreResponse.Id)

	var writeAuthorizationModelRequest client.ClientWriteAuthorizationModelRequest
	err = json.Unmarshal([]byte(authModel), &writeAuthorizationModelRequest)
	s.Require().NoError(err)

	writeAuthorizationModelResponse, err := s.fga.WriteAuthorizationModel(context.Background()).Body(writeAuthorizationModelRequest).Execute()
	s.Require().NoError(err)

	s.authModelID = writeAuthorizationModelResponse.AuthorizationModelId

	clientWriteResponse, err := s.fga.WriteTuples(context.Background()).Options(client.ClientWriteOptions{AuthorizationModelId: s.authModelID}).Body(client.ClientWriteTuplesBody{
		{
			User:     "user:*",
			Relation: "user",
			Object:   "server:lxd",
		},
		{
			User:     "server:lxd",
			Relation: "server",
			Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
		},
		{
			User:     "server:lxd",
			Relation: "server",
			Object:   "cluster_member:node01",
		},
		{
			User:     "server:lxd",
			Relation: "server",
			Object:   "cluster_group:group01",
		},
		{
			User:     "server:lxd",
			Relation: "server",
			Object:   "storage_pool:pool01",
		},
		{
			User:     "server:lxd",
			Relation: "server",
			Object:   "project:project01",
		},
		{
			User:     "project:project01",
			Relation: "project",
			Object:   "image:image01",
		},
		{
			User:     "project:project01",
			Relation: "project",
			Object:   "instance:instance01",
		},
		{
			User:     "project:project01",
			Relation: "project",
			Object:   "network:network01",
		},
		{
			User:     "project:project01",
			Relation: "project",
			Object:   "network_acl:network_acl01",
		},
		{
			User:     "project:project01",
			Relation: "project",
			Object:   "network_zone:network_zone01",
		},
		{
			User:     "project:project01",
			Relation: "project",
			Object:   "network_forward:network_forward01",
		},
		{
			User:     "project:project01",
			Relation: "project",
			Object:   "network_load_balancer:network_load_balancer01",
		},
		{
			User:     "project:project01",
			Relation: "project",
			Object:   "network_peer:network_peer01",
		},
		{
			User:     "project:project01",
			Relation: "project",
			Object:   "profile:profile01",
		},
		{
			User:     "project:project01",
			Relation: "project",
			Object:   "storage_pool_volume:storage_pool_volume01",
		},
		{
			User:     "project:project01",
			Relation: "project",
			Object:   "storage_bucket:storage_bucket01",
		},
	}).Execute()
	s.Require().NoError(err)

	s.Require().Len(clientWriteResponse.Deletes, 0)
	for _, write := range clientWriteResponse.Writes {
		s.Require().NoError(write.Error)
	}
}

func (s *openFGASuite) TearDownTest() {
	_, err := s.fga.DeleteStore(context.Background()).Execute()
	s.Require().NoError(err)
}

func (s *openFGASuite) TestPublic() {
	type test struct {
		description string
		allowed     bool
		request     client.ClientCheckRequest
	}

	tests := []test{
		{
			description: "User with no relations should not be able to edit the server",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "User with no relations should be able to view the server",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "User with no relations should not be able to create a storage pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_storage_pool",
				Object:   "server:lxd",
			},
		},
		{
			description: "User with no relations should not be able to create a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_project",
				Object:   "server:lxd",
			},
		},
		{
			description: "User with no relations should not be able to view server resources",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view_resources",
				Object:   "server:lxd",
			},
		},
		{
			description: "User with no relations should not be able to create a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_certificate",
				Object:   "server:lxd",
			},
		},
		{
			description: "User with no relations should not be able to edit cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "User with no relations should not be able to view cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "User with no relations should not be able to add a new member to a cluster",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_cluster_member",
				Object:   "server:lxd",
			},
		},
		{
			description: "User with no relations should not be able to create a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_cluster_group",
				Object:   "server:lxd",
			},
		},
		{
			description: "User with no relations should not be able to view server metrics",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view_metrics",
				Object:   "server:lxd",
			},
		},
		{
			description: "User with no relations should not be able to edit a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "User with no relations should not be able to view a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "User with no relations should not be able to edit cluster member config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "User with no relations should not be able to view a cluster member",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "User with no relations should not be able to edit a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "User with no relations should not be able to view a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "User with no relations should not be able to edit a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "User with no relations should not be able to view a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "User with no relations should not be able to edit a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "project:project01",
			},
		},
		{
			description: "User with no relations should not be able to view a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "project:project01",
			},
		},
		{
			description: "User with no relations should not be able to create instances in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_instances",
				Object:   "project:project01",
			},
		},
		{
			description: "User with no relations should not be able to create images in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_images",
				Object:   "project:project01",
			},
		},
		{
			description: "User with no relations should not be able to create networks in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "User with no relations should not be able to create network ACLs in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_network_acls",
				Object:   "project:project01",
			},
		},
		{
			description: "User with no relations should not be able to create network zones in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_network_zones",
				Object:   "project:project01",
			},
		},
		{
			description: "User with no relations should not be able to create network forwards in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_network_forwards",
				Object:   "project:project01",
			},
		},
		{
			description: "User with no relations should not be able to create network load balancers in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "User with no relations should not be able to create network peers in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_network_peers",
				Object:   "project:project01",
			},
		},
		{
			description: "User with no relations should not be able to create profiles in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_profiles",
				Object:   "project:project01",
			},
		},
		{
			description: "User with no relations should not be able to create storage pool volumes in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_storage_pool_volumes",
				Object:   "project:project01",
			},
		},
		{
			description: "User with no relations should not be able to create storage_buckets in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "User with no relations should not be able to edit an image",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "image:image01",
			},
		},
		{
			description: "User with no relations should not be able to view an image",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "image:image01",
			},
		},
		{
			description: "User with no relations should not be able to edit an instance",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User with no relations should not be able to view an instance",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User with no relations should not be able to update and instances' state",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User with no relations should not be able to manage an instances' snapshots",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_manage_snapshots",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User with no relations should not be able to manage an instances' backups",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User with no relations should not be able to connect to an instance via sftp",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_connect_sftp",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User with no relations should not be able to push/pull files into an instance",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_access_files",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User with no relations should not be able to access an instances' console",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_access_console",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User with no relations should not be able to exec in an instance",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_exec",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User with no relations should not be able to edit a network",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "network:network01",
			},
		},
		{
			description: "User with no relations should not be able to view a network",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "network:network01",
			},
		},
		{
			description: "User with no relations should not be able to edit a network ACL",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "User with no relations should not be able to view a network ACL",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "User with no relations should not be able to edit a network forward",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "User with no relations should not be able to view a network forward",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "User with no relations should not be able to edit a network load balancer",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "User with no relations should not be able to view a network load balancer",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "User with no relations should not be able to edit a network peer",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "User with no relations should not be able to view a network peer",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "User with no relations should not be able to edit a profile",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "profile:profile01",
			},
		},
		{
			description: "User with no relations should not be able to view a profile",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "profile:profile01",
			},
		},
		{
			description: "User with no relations should not be able to edit a storage pool volume",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "User with no relations should not be able to view a storage pool volume",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "User with no relations should not be able to edit a storage bucket",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_edit",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
		{
			description: "User with no relations should not be able to view a storage bucket",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:anyone",
				Relation: "can_view",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
	}

	for i, test := range tests {
		s.T().Logf("Case %d: %s", i, test.description)

		checkResponse, err := s.fga.Check(context.Background()).Options(client.ClientCheckOptions{AuthorizationModelId: s.authModelID}).Body(test.request).Execute()
		s.Require().NoError(err)
		s.Equal(test.allowed, checkResponse.GetAllowed())
	}
}

func (s *openFGASuite) TestServerAdmin() {
	clientWriteResponse, err := s.fga.WriteTuples(context.Background()).Options(client.ClientWriteOptions{AuthorizationModelId: s.authModelID}).Body(client.ClientWriteTuplesBody{
		{
			User:     "group:server_admins#member",
			Relation: "admin",
			Object:   "server:lxd",
		},
		{
			User:     "user:server_admin",
			Relation: "member",
			Object:   "group:server_admins",
		},
	}).Execute()
	s.Require().NoError(err)

	s.Require().Len(clientWriteResponse.Deletes, 0)
	for _, write := range clientWriteResponse.Writes {
		s.Require().NoError(write.Error)
	}

	tests := []struct {
		description string
		allowed     bool
		request     client.ClientCheckRequest
	}{
		{
			description: "Server admin should be able to edit the server",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server admin should be able to view the server",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server admin should be able to create a storage pool",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_storage_pool",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server admin should be able to create a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_project",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server admin should be able to view server resources",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view_resources",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server admin should be able to create a certificate",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_certificate",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server admin should be able to edit cluster config",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server admin should be able to view cluster config",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server admin should be able to add a new member to a cluster",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_cluster_member",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server admin should be able to create a cluster group",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_cluster_group",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server admin should be able to view server metrics",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view_metrics",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server admin should be able to edit a certificate",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Server admin should be able to view a certificate",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Server admin should be able to edit cluster member config",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Server admin should be able to view a cluster member",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Server admin should be able to edit a cluster group",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Server admin should be able to view a cluster group",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Server admin should be able to edit a storage_pool",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Server admin should be able to view a storage_pool",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Server admin should be able to edit a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "project:project01",
			},
		},
		{
			description: "Server admin should be able to view a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "project:project01",
			},
		},
		{
			description: "Server admin should be able to create instances in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_instances",
				Object:   "project:project01",
			},
		},
		{
			description: "Server admin should be able to create images in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_images",
				Object:   "project:project01",
			},
		},
		{
			description: "Server admin should be able to create networks in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Server admin should be able to create network ACLs in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_network_acls",
				Object:   "project:project01",
			},
		},
		{
			description: "Server admin should be able to create network zones in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_network_zones",
				Object:   "project:project01",
			},
		},
		{
			description: "Server admin should be able to create network forwards in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_network_forwards",
				Object:   "project:project01",
			},
		},
		{
			description: "Server admin should be able to create network load balancers in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Server admin should be able to create network peers in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_network_peers",
				Object:   "project:project01",
			},
		},
		{
			description: "Server admin should be able to create profiles in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_profiles",
				Object:   "project:project01",
			},
		},
		{
			description: "Server admin should be able to create storage pool volumes in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_storage_pool_volumes",
				Object:   "project:project01",
			},
		},
		{
			description: "Server admin should be able to create storage_buckets in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Server admin should be able to edit an image",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "image:image01",
			},
		},
		{
			description: "Server admin should be able to view an image",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "image:image01",
			},
		},
		{
			description: "Server admin should be able to edit an instance",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server admin should be able to view an instance",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server admin should be able to update and instances' state",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server admin should be able to manage an instances' snapshots",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_manage_snapshots",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server admin should be able to manage an instances' backups",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server admin should be able to connect to an instance via sftp",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_connect_sftp",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server admin should be able to push/pull files into an instance",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_access_files",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server admin should be able to access an instances' console",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_access_console",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server admin should be able to exec in an instance",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_exec",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server admin should be able to edit a network",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "network:network01",
			},
		},
		{
			description: "Server admin should be able to view a network",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "network:network01",
			},
		},
		{
			description: "Server admin should be able to edit a network ACL",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Server admin should be able to view a network ACL",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Server admin should be able to edit a network forward",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Server admin should be able to view a network forward",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Server admin should be able to edit a network load balancer",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Server admin should be able to view a network load balancer",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Server admin should be able to edit a network peer",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Server admin should be able to view a network peer",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Server admin should be able to edit a profile",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Server admin should be able to view a profile",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Server admin should be able to edit a storage pool volume",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Server admin should be able to view a storage pool volume",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Server admin should be able to edit a storage bucket",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_edit",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
		{
			description: "Server admin should be able to view a storage bucket",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_admin",
				Relation: "can_view",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
	}

	for i, test := range tests {
		s.T().Logf("Case %d: %s", i, test.description)

		checkResponse, err := s.fga.Check(context.Background()).Options(client.ClientCheckOptions{AuthorizationModelId: s.authModelID}).Body(test.request).Execute()
		s.Require().NoError(err)
		s.Equal(test.allowed, checkResponse.GetAllowed())
	}
}

func (s *openFGASuite) TestServerOperator() {
	clientWriteResponse, err := s.fga.WriteTuples(context.Background()).Options(client.ClientWriteOptions{AuthorizationModelId: s.authModelID}).Body(client.ClientWriteTuplesBody{
		{
			User:     "group:server_operators#member",
			Relation: "operator",
			Object:   "server:lxd",
		},
		{
			User:     "user:server_operator",
			Relation: "member",
			Object:   "group:server_operators",
		},
	}).Execute()
	s.Require().NoError(err)

	s.Require().Len(clientWriteResponse.Deletes, 0)
	for _, write := range clientWriteResponse.Writes {
		s.Require().NoError(write.Error)
	}

	type test struct {
		description string
		allowed     bool
		request     client.ClientCheckRequest
	}
	tests := []test{
		{
			description: "Server operator should not be able to edit the server",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server operator should be able to view the server",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server operator should not be able to create a storage pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_storage_pool",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server operator should be able to create a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_project",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server operator should be able to view server resources",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view_resources",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server operator should not be able to create a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_certificate",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server operator should not be able to edit cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server operator should be able to view cluster config",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server operator should not be able to add a new member to a cluster",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_cluster_member",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server operator should not be able to create a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_cluster_group",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server operator should be able to view server metrics",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view_metrics",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server operator should not be able to edit a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Server operator should be able to view a certificate",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Server operator should not be able to edit cluster member config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Server operator should be able to view a cluster member",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Server operator should not be able to edit a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Server operator should be able to view a cluster group",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Server operator should not be able to edit a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Server operator should be able to view a storage_pool",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Server operator should be able to edit a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "project:project01",
			},
		},
		{
			description: "Server operator should be able to view a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "project:project01",
			},
		},
		{
			description: "Server operator should be able to create instances in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_instances",
				Object:   "project:project01",
			},
		},
		{
			description: "Server operator should be able to create images in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_images",
				Object:   "project:project01",
			},
		},
		{
			description: "Server operator should be able to create networks in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Server operator should be able to create network ACLs in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_network_acls",
				Object:   "project:project01",
			},
		},
		{
			description: "Server operator should be able to create network zones in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_network_zones",
				Object:   "project:project01",
			},
		},
		{
			description: "Server operator should be able to create network forwards in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_network_forwards",
				Object:   "project:project01",
			},
		},
		{
			description: "Server operator should be able to create network load balancers in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Server operator should be able to create network peers in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_network_peers",
				Object:   "project:project01",
			},
		},
		{
			description: "Server operator should be able to create profiles in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_profiles",
				Object:   "project:project01",
			},
		},
		{
			description: "Server operator should be able to create storage pool volumes in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_storage_pool_volumes",
				Object:   "project:project01",
			},
		},
		{
			description: "Server operator should be able to create storage_buckets in a project",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Server operator should be able to edit an image",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "image:image01",
			},
		},
		{
			description: "Server operator should be able to view an image",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "image:image01",
			},
		},
		{
			description: "Server operator should be able to edit an instance",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server operator should be able to view an instance",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server operator should be able to update and instances' state",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server operator should be able to manage an instances' snapshots",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_manage_snapshots",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server operator should be able to manage an instances' backups",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server operator should be able to connect to an instance via sftp",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_connect_sftp",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server operator should be able to push/pull files into an instance",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_access_files",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server operator should be able to access an instances' console",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_access_console",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server operator should be able to exec in an instance",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_exec",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server operator should be able to edit a network",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "network:network01",
			},
		},
		{
			description: "Server operator should be able to view a network",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "network:network01",
			},
		},
		{
			description: "Server operator should be able to edit a network ACL",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Server operator should be able to view a network ACL",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Server operator should be able to edit a network forward",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Server operator should be able to view a network forward",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Server operator should be able to edit a network load balancer",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Server operator should be able to view a network load balancer",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Server operator should be able to edit a network peer",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Server operator should be able to view a network peer",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Server operator should be able to edit a profile",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Server operator should be able to view a profile",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Server operator should be able to edit a storage pool volume",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Server operator should be able to view a storage pool volume",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Server operator should be able to edit a storage bucket",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_edit",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
		{
			description: "Server operator should be able to view a storage bucket",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_operator",
				Relation: "can_view",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
	}

	for i, test := range tests {
		s.T().Logf("Case %d: %s", i, test.description)

		checkResponse, err := s.fga.Check(context.Background()).Options(client.ClientCheckOptions{AuthorizationModelId: s.authModelID}).Body(test.request).Execute()
		s.Require().NoError(err)
		s.Equal(test.allowed, checkResponse.GetAllowed())
	}
}

func (s *openFGASuite) TestServerViewer() {
	clientWriteResponse, err := s.fga.WriteTuples(context.Background()).Options(client.ClientWriteOptions{AuthorizationModelId: s.authModelID}).Body(client.ClientWriteTuplesBody{
		{
			User:     "group:server_viewers#member",
			Relation: "viewer",
			Object:   "server:lxd",
		},
		{
			User:     "user:server_viewer",
			Relation: "member",
			Object:   "group:server_viewers",
		},
	}).Execute()
	s.Require().NoError(err)

	s.Require().Len(clientWriteResponse.Deletes, 0)
	for _, write := range clientWriteResponse.Writes {
		s.Require().NoError(write.Error)
	}

	type test struct {
		description string
		allowed     bool
		request     client.ClientCheckRequest
	}

	tests := []test{
		{
			description: "Server viewer should not be able to edit the server",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server viewer should be able to view the server",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server viewer should not be able to create a storage pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_storage_pool",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server viewer should not be able to create a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_project",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server viewer should be able to view server resources",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view_resources",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server viewer should not be able to create a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_certificate",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server viewer should not be able to edit cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server viewer should be able to view cluster config",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server viewer should not be able to add a new member to a cluster",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_cluster_member",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server viewer should not be able to create a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_cluster_group",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server viewer should be able to view server metrics",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view_metrics",
				Object:   "server:lxd",
			},
		},
		{
			description: "Server viewer should not be able to edit a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Server viewer should be able to view a certificate",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Server viewer should not be able to edit cluster member config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Server viewer should be able to view a cluster member",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Server viewer should not be able to edit a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Server viewer should be able to view a cluster group",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Server viewer should not be able to edit a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Server viewer should be able to view a storage_pool",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Server viewer should not be able to edit a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "project:project01",
			},
		},
		{
			description: "Server viewer should not be able to view a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "project:project01",
			},
		},
		{
			description: "Server viewer should not be able to create instances in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_instances",
				Object:   "project:project01",
			},
		},
		{
			description: "Server viewer should not be able to create images in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_images",
				Object:   "project:project01",
			},
		},
		{
			description: "Server viewer should not be able to create networks in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Server viewer should not be able to create network ACLs in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_network_acls",
				Object:   "project:project01",
			},
		},
		{
			description: "Server viewer should not be able to create network zones in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_network_zones",
				Object:   "project:project01",
			},
		},
		{
			description: "Server viewer should not be able to create network forwards in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_network_forwards",
				Object:   "project:project01",
			},
		},
		{
			description: "Server viewer should not be able to create network load balancers in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Server viewer should not be able to create network peers in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_network_peers",
				Object:   "project:project01",
			},
		},
		{
			description: "Server viewer should not be able to create profiles in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_profiles",
				Object:   "project:project01",
			},
		},
		{
			description: "Server viewer should not be able to create storage pool volumes in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_storage_pool_volumes",
				Object:   "project:project01",
			},
		},
		{
			description: "Server viewer should not be able to create storage_buckets in a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Server viewer should not be able to edit an image",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "image:image01",
			},
		},
		{
			description: "Server viewer should not be able to view an image",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "image:image01",
			},
		},
		{
			description: "Server viewer should not be able to edit an instance",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server viewer should not be able to view an instance",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server viewer should not be able to update and instances' state",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server viewer should not be able to manage an instances' snapshots",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_manage_snapshots",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server viewer should not be able to manage an instances' backups",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server viewer should not be able to connect to an instance via sftp",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_connect_sftp",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server viewer should not be able to push/pull files into an instance",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_access_files",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server viewer should not be able to access an instances' console",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_access_console",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server viewer should not be able to exec in an instance",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_exec",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Server viewer should not be able to edit a network",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "network:network01",
			},
		},
		{
			description: "Server viewer should not be able to view a network",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "network:network01",
			},
		},
		{
			description: "Server viewer should not be able to edit a network ACL",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Server viewer should not be able to view a network ACL",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Server viewer should not be able to edit a network forward",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Server viewer should not be able to view a network forward",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Server viewer should not be able to edit a network load balancer",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Server viewer should not be able to view a network load balancer",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Server viewer should not be able to edit a network peer",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Server viewer should not be able to view a network peer",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Server viewer should not be able to edit a profile",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Server viewer should not be able to view a profile",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Server viewer should not be able to edit a storage pool volume",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Server viewer should not be able to view a storage pool volume",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Server viewer should not be able to edit a storage bucket",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_edit",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
		{
			description: "Server viewer should not be able to view a storage bucket",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:server_viewer",
				Relation: "can_view",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
	}

	for i, test := range tests {
		s.T().Logf("Case %d: %s", i, test.description)

		checkResponse, err := s.fga.Check(context.Background()).Options(client.ClientCheckOptions{AuthorizationModelId: s.authModelID}).Body(test.request).Execute()
		s.Require().NoError(err)
		s.Equal(test.allowed, checkResponse.GetAllowed())
	}
}

func (s *openFGASuite) TestProjectManager() {
	clientWriteResponse, err := s.fga.WriteTuples(context.Background()).Options(client.ClientWriteOptions{AuthorizationModelId: s.authModelID}).Body(client.ClientWriteTuplesBody{
		{
			User:     "group:project01_managers#member",
			Relation: "manager",
			Object:   "project:project01",
		},
		{
			// A manager of a project needs view permissions on the server to allow targeting of cluster groups
			// or to be able to view available storage pools for creating volumes etc.
			User:     "group:project01_managers#member",
			Relation: "viewer",
			Object:   "server:lxd",
		},
		{
			User:     "user:project01_manager",
			Relation: "member",
			Object:   "group:project01_managers",
		},
	}).Execute()
	s.Require().NoError(err)

	s.Require().Len(clientWriteResponse.Deletes, 0)
	for _, write := range clientWriteResponse.Writes {
		s.Require().NoError(write.Error)
	}

	type test struct {
		description string
		allowed     bool
		request     client.ClientCheckRequest
	}

	tests := []test{
		{
			description: "Manager of project01 should not be able to edit the server",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of project01 should be able to view the server",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of project01 should not be able to create a storage pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_storage_pool",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of project01 should not be able to create a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_project",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of project01 should be able to view server resources",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view_resources",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of project01 should not be able to create a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_certificate",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of project01 should not be able to edit cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of project01 should be able to view cluster config",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of project01 should not be able to add a new member to a cluster",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_cluster_member",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of project01 should not be able to create a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_cluster_group",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of project01 should be able to view server metrics",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view_metrics",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of project01 should not be able to edit a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Manager of project01 should be able to view a certificate",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Manager of project01 should not be able to edit cluster member config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Manager of project01 should be able to view a cluster member",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Manager of project01 should not be able to edit a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Manager of project01 should be able to view a cluster group",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Manager of project01 should not be able to edit a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Manager of project01 should be able to view a storage_pool",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Manager of project01 should be able to edit project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of project01 should be able to view project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of project01 should be able to create instances in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_instances",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of project01 should be able to create images in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_images",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of project01 should be able to create networks in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of project01 should be able to create network ACLs in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_network_acls",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of project01 should be able to create network zones in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_network_zones",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of project01 should be able to create network forwards in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_network_forwards",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of project01 should be able to create network load balancers in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of project01 should be able to create network peers in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_network_peers",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of project01 should be able to create profiles in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_profiles",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of project01 should be able to create storage pool volumes in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_storage_pool_volumes",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of project01 should be able to create storage_buckets in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of project01 should be able to edit an image in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "image:image01",
			},
		},
		{
			description: "Manager of project01 should be able to view an image in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "image:image01",
			},
		},
		{
			description: "Manager of project01 should be able to edit an instance in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of project01 should be able to view an instance in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of project01 should be able to update and instances' state in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of project01 should be able to manage an instances' snapshots in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_manage_snapshots",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of project01 should be able to manage an instances' backups in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of project01 should be able to connect to an instance via sftp in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_connect_sftp",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of project01 should be able to push/pull files into an instance in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_access_files",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of project01 should be able to access an instances' console in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_access_console",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of project01 should be able to exec in an instance in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_exec",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of project01 should be able to edit a network in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "network:network01",
			},
		},
		{
			description: "Manager of project01 should be able to view a network in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "network:network01",
			},
		},
		{
			description: "Manager of project01 should be able to edit a network ACL in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Manager of project01 should be able to view a network ACL in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Manager of project01 should be able to edit a network forward in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Manager of project01 should be able to view a network forward in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Manager of project01 should be able to edit a network load balancer in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Manager of project01 should be able to view a network load balancer in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Manager of project01 should be able to edit a network peer in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Manager of project01 should be able to view a network peer in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Manager of project01 should be able to edit a profile in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Manager of project01 should be able to view a profile in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Manager of project01 should be able to edit a storage pool volume in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Manager of project01 should be able to view a storage pool volume in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Manager of project01 should be able to edit a storage bucket in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_edit",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
		{
			description: "Manager of project01 should be able to view a storage bucket in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_manager",
				Relation: "can_view",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
	}

	for i, test := range tests {
		s.T().Logf("Case %d: %s", i, test.description)

		checkResponse, err := s.fga.Check(context.Background()).Options(client.ClientCheckOptions{AuthorizationModelId: s.authModelID}).Body(test.request).Execute()
		s.Require().NoError(err)
		s.Equal(test.allowed, checkResponse.GetAllowed())
	}
}

func (s *openFGASuite) TestProjectOperator() {
	clientWriteResponse, err := s.fga.WriteTuples(context.Background()).Options(client.ClientWriteOptions{AuthorizationModelId: s.authModelID}).Body(client.ClientWriteTuplesBody{
		{
			User:     "group:project01_operators#member",
			Relation: "operator",
			Object:   "project:project01",
		},
		{
			// An operator of a project needs view permissions on the server to allow targeting of cluster groups
			// or to be able to view available storage pools for creating volumes etc.
			User:     "group:project01_operators#member",
			Relation: "viewer",
			Object:   "server:lxd",
		},
		{
			User:     "user:project01_operator",
			Relation: "member",
			Object:   "group:project01_operators",
		},
	}).Execute()
	s.Require().NoError(err)

	s.Require().Len(clientWriteResponse.Deletes, 0)
	for _, write := range clientWriteResponse.Writes {
		s.Require().NoError(write.Error)
	}

	type test struct {
		description string
		allowed     bool
		request     client.ClientCheckRequest
	}

	tests := []test{
		{
			description: "Operator of project01 should not be able to edit the server",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of project01 should be able to view the server",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of project01 should not be able to create a storage pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_storage_pool",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of project01 should not be able to create a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_project",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of project01 should be able to view server resources",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view_resources",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of project01 should not be able to create a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_certificate",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of project01 should not be able to edit cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of project01 should be able to view cluster config",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of project01 should not be able to add a new member to a cluster",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_cluster_member",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of project01 should not be able to create a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_cluster_group",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of project01 should be able to view server metrics",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view_metrics",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of project01 should not be able to edit a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Operator of project01 should be able to view a certificate",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Operator of project01 should not be able to edit cluster member config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Operator of project01 should be able to view a cluster member",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Operator of project01 should not be able to edit a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Operator of project01 should be able to view a cluster group",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Operator of project01 should not be able to edit a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Operator of project01 should be able to view a storage_pool",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Operator of project01 should not be able to edit project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of project01 should be able to view project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of project01 should be able to create instances in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_instances",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of project01 should be able to create images in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_images",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of project01 should be able to create networks in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of project01 should be able to create network ACLs in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_network_acls",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of project01 should be able to create network zones in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_network_zones",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of project01 should be able to create network forwards in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_network_forwards",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of project01 should be able to create network load balancers in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of project01 should be able to create network peers in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_network_peers",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of project01 should be able to create profiles in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_profiles",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of project01 should be able to create storage pool volumes in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_storage_pool_volumes",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of project01 should be able to create storage_buckets in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of project01 should be able to edit an image in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "image:image01",
			},
		},
		{
			description: "Operator of project01 should be able to view an image in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "image:image01",
			},
		},
		{
			description: "Operator of project01 should be able to edit an instance in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of project01 should be able to view an instance in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of project01 should be able to update and instances' state in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of project01 should be able to manage an instances' snapshots in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_manage_snapshots",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of project01 should be able to manage an instances' backups in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of project01 should be able to connect to an instance via sftp in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_connect_sftp",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of project01 should be able to push/pull files into an instance in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_access_files",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of project01 should be able to access an instances' console in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_access_console",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of project01 should be able to exec in an instance in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_exec",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of project01 should be able to edit a network in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "network:network01",
			},
		},
		{
			description: "Operator of project01 should be able to view a network in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "network:network01",
			},
		},
		{
			description: "Operator of project01 should be able to edit a network ACL in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Operator of project01 should be able to view a network ACL in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Operator of project01 should be able to edit a network forward in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Operator of project01 should be able to view a network forward in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Operator of project01 should be able to edit a network load balancer in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Operator of project01 should be able to view a network load balancer in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Operator of project01 should be able to edit a network peer in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Operator of project01 should be able to view a network peer in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Operator of project01 should be able to edit a profile in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Operator of project01 should be able to view a profile in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Operator of project01 should be able to edit a storage pool volume in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Operator of project01 should be able to view a storage pool volume in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Operator of project01 should be able to edit a storage bucket in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_edit",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
		{
			description: "Operator of project01 should be able to view a storage bucket in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_operator",
				Relation: "can_view",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
	}

	for i, test := range tests {
		s.T().Logf("Case %d: %s", i, test.description)

		checkResponse, err := s.fga.Check(context.Background()).Options(client.ClientCheckOptions{AuthorizationModelId: s.authModelID}).Body(test.request).Execute()
		s.Require().NoError(err)
		s.Equal(test.allowed, checkResponse.GetAllowed())
	}
}

func (s *openFGASuite) TestProjectViewer() {
	clientWriteResponse, err := s.fga.WriteTuples(context.Background()).Options(client.ClientWriteOptions{AuthorizationModelId: s.authModelID}).Body(client.ClientWriteTuplesBody{
		{
			User:     "group:project01_viewers#member",
			Relation: "viewer",
			Object:   "project:project01",
		},
		{
			User:     "user:project01_viewer",
			Relation: "member",
			Object:   "group:project01_viewers",
		},
	}).Execute()
	s.Require().NoError(err)

	s.Require().Len(clientWriteResponse.Deletes, 0)
	for _, write := range clientWriteResponse.Writes {
		s.Require().NoError(write.Error)
	}

	type test struct {
		description string
		allowed     bool
		request     client.ClientCheckRequest
	}

	tests := []test{
		{
			description: "Viewer of project01 should not be able to edit the server",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Viewer of project01 should be able to view the server",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Viewer of project01 should not be able to create a storage pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_storage_pool",
				Object:   "server:lxd",
			},
		},
		{
			description: "Viewer of project01 should not be able to create a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_project",
				Object:   "server:lxd",
			},
		},
		{
			description: "Viewer of project01 should not be able to view server resources",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view_resources",
				Object:   "server:lxd",
			},
		},
		{
			description: "Viewer of project01 should not be able to create a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_certificate",
				Object:   "server:lxd",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Viewer of project01 should not be able to view cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Viewer of project01 should not be able to add a new member to a cluster",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_cluster_member",
				Object:   "server:lxd",
			},
		},
		{
			description: "Viewer of project01 should not be able to create a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_cluster_group",
				Object:   "server:lxd",
			},
		},
		{
			description: "Viewer of project01 should not be able to view server metrics",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view_metrics",
				Object:   "server:lxd",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Viewer of project01 should not be able to view a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit cluster member config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Viewer of project01 should not be able to view a cluster member",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Viewer of project01 should not be able to view a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Viewer of project01 should not be able to view a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "project:project01",
			},
		},
		{
			description: "Viewer of project01 should be able to view project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "project:project01",
			},
		},
		{
			description: "Viewer of project01 should not be able to create instances in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_instances",
				Object:   "project:project01",
			},
		},
		{
			description: "Viewer of project01 should not be able to create images in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_images",
				Object:   "project:project01",
			},
		},
		{
			description: "Viewer of project01 should not be able to create networks in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Viewer of project01 should not be able to create network ACLs in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_network_acls",
				Object:   "project:project01",
			},
		},
		{
			description: "Viewer of project01 should not be able to create network zones in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_network_zones",
				Object:   "project:project01",
			},
		},
		{
			description: "Viewer of project01 should not be able to create network forwards in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_network_forwards",
				Object:   "project:project01",
			},
		},
		{
			description: "Viewer of project01 should not be able to create network load balancers in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Viewer of project01 should not be able to create network peers in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_network_peers",
				Object:   "project:project01",
			},
		},
		{
			description: "Viewer of project01 should not be able to create profiles in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_profiles",
				Object:   "project:project01",
			},
		},
		{
			description: "Viewer of project01 should not be able to create storage pool volumes in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_storage_pool_volumes",
				Object:   "project:project01",
			},
		},
		{
			description: "Viewer of project01 should not be able to create storage_buckets in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit an image in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "image:image01",
			},
		},
		{
			description: "Viewer of project01 should be able to view an image in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "image:image01",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit an instance in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Viewer of project01 should be able to view an instance in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Viewer of project01 should not be able to update and instances' state in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Viewer of project01 should not be able to manage an instances' snapshots in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_manage_snapshots",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Viewer of project01 should not be able to manage an instances' backups in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Viewer of project01 should not be able to connect to an instance via sftp in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_connect_sftp",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Viewer of project01 should not be able to push/pull files into an instance in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_access_files",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Viewer of project01 should not be able to access an instances' console in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_access_console",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Viewer of project01 should not be able to exec in an instance in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_exec",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit a network in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "network:network01",
			},
		},
		{
			description: "Viewer of project01 should be able to view a network in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "network:network01",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit a network ACL in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Viewer of project01 should be able to view a network ACL in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit a network forward in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Viewer of project01 should be able to view a network forward in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit a network load balancer in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Viewer of project01 should be able to view a network load balancer in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit a network peer in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Viewer of project01 should be able to view a network peer in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit a profile in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Viewer of project01 should be able to view a profile in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit a storage pool volume in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Viewer of project01 should be able to view a storage pool volume in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Viewer of project01 should not be able to edit a storage bucket in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_edit",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
		{
			description: "Viewer of project01 should be able to view a storage bucket in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:project01_viewer",
				Relation: "can_view",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
	}

	for i, test := range tests {
		s.T().Logf("Case %d: %s", i, test.description)

		checkResponse, err := s.fga.Check(context.Background()).Options(client.ClientCheckOptions{AuthorizationModelId: s.authModelID}).Body(test.request).Execute()
		s.Require().NoError(err)
		s.Equal(test.allowed, checkResponse.GetAllowed())
	}
}

func (s *openFGASuite) TestInstanceManager() {
	clientWriteResponse, err := s.fga.WriteTuples(context.Background()).Options(client.ClientWriteOptions{AuthorizationModelId: s.authModelID}).Body(client.ClientWriteTuplesBody{
		{
			User:     "group:instance01_managers#member",
			Relation: "manager",
			Object:   "instance:instance01",
		},
		{
			// A manager of an instance will need view permissions on a project to see e.g. which storage volumes are available to be mounted.
			User:     "group:instance01_managers#member",
			Relation: "viewer",
			Object:   "project:project01",
		},
		{
			User:     "user:instance01_manager",
			Relation: "member",
			Object:   "group:instance01_managers",
		},
	}).Execute()
	s.Require().NoError(err)

	s.Require().Len(clientWriteResponse.Deletes, 0)
	for _, write := range clientWriteResponse.Writes {
		s.Require().NoError(write.Error)
	}

	type test struct {
		description string
		allowed     bool
		request     client.ClientCheckRequest
	}

	tests := []test{
		{
			description: "Manager of instance01 should not be able to edit the server",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of instance01 should be able to view the server",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of instance01 should not be able to create a storage pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_storage_pool",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of instance01 should not be able to create a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_project",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of instance01 should not be able to view server resources",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view_resources",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of instance01 should not be able to create a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_certificate",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of instance01 should not be able to view cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of instance01 should not be able to add a new member to a cluster",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_cluster_member",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of instance01 should not be able to create a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_cluster_group",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of instance01 should not be able to view server metrics",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view_metrics",
				Object:   "server:lxd",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Manager of instance01 should not be able to view a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit cluster member config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Manager of instance01 should not be able to view a cluster member",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Manager of instance01 should not be able to view a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Manager of instance01 should not be able to view a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of instance01 should be able to view project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of instance01 should not be able to create instances in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_instances",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of instance01 should not be able to create images in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_images",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of instance01 should not be able to create networks in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of instance01 should not be able to create network ACLs in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_network_acls",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of instance01 should not be able to create network zones in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_network_zones",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of instance01 should not be able to create network forwards in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_network_forwards",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of instance01 should not be able to create network load balancers in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of instance01 should not be able to create network peers in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_network_peers",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of instance01 should not be able to create profiles in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_profiles",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of instance01 should not be able to create storage pool volumes in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_storage_pool_volumes",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of instance01 should not be able to create storage_buckets in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit an image in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "image:image01",
			},
		},
		{
			description: "Manager of instance01 should be able to view an image in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "image:image01",
			},
		},
		{
			description: "Manager of instance01 should be able to edit instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of instance01 should be able to view instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of instance01 should be able to update the state of instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of instance01 should be able to manage snapshots of instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_manage_snapshots",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of instance01 should be able to manage backups of instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of instance01 should be able to connect to instance01 via sftp in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_connect_sftp",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of instance01 should be able to push/pull files into instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_access_files",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of instance01 should be able to access instance01's console",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_access_console",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of instance01 should be able to exec into instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_exec",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit a network in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "network:network01",
			},
		},
		{
			description: "Manager of instance01 should be able to view a network in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "network:network01",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit a network ACL in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Manager of instance01 should be able to view a network ACL in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit a network forward in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Manager of instance01 should be able to view a network forward in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit a network load balancer in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Manager of instance01 should be able to view a network load balancer in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit a network peer in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Manager of instance01 should be able to view a network peer in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit a profile in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Manager of instance01 should be able to view a profile in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit a storage pool volume in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Manager of instance01 should be able to view a storage pool volume in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Manager of instance01 should not be able to edit a storage bucket in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_edit",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
		{
			description: "Manager of instance01 should be able to view a storage bucket in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_manager",
				Relation: "can_view",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
	}

	for i, test := range tests {
		s.T().Logf("Case %d: %s", i, test.description)

		checkResponse, err := s.fga.Check(context.Background()).Options(client.ClientCheckOptions{AuthorizationModelId: s.authModelID}).Body(test.request).Execute()
		s.Require().NoError(err)
		s.Equal(test.allowed, checkResponse.GetAllowed())
	}
}

func (s *openFGASuite) TestInstanceOperator() {
	clientWriteResponse, err := s.fga.WriteTuples(context.Background()).Options(client.ClientWriteOptions{AuthorizationModelId: s.authModelID}).Body(client.ClientWriteTuplesBody{
		{
			User:     "group:instance01_operators#member",
			Relation: "operator",
			Object:   "instance:instance01",
		},
		{
			User:     "user:instance01_operator",
			Relation: "member",
			Object:   "group:instance01_operators",
		},
	}).Execute()
	s.Require().NoError(err)

	s.Require().Len(clientWriteResponse.Deletes, 0)
	for _, write := range clientWriteResponse.Writes {
		s.Require().NoError(write.Error)
	}

	type test struct {
		description string
		allowed     bool
		request     client.ClientCheckRequest
	}

	tests := []test{
		{
			description: "Operator of instance01 should not be able to edit the server",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of instance01 should be able to view the server",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of instance01 should not be able to create a storage pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_storage_pool",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of instance01 should not be able to create a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_project",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of instance01 should not be able to view server resources",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view_resources",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of instance01 should not be able to create a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_certificate",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of instance01 should not be able to view cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of instance01 should not be able to add a new member to a cluster",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_cluster_member",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of instance01 should not be able to create a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_cluster_group",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of instance01 should not be able to view server metrics",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view_metrics",
				Object:   "server:lxd",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Operator of instance01 should not be able to view a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit cluster member config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Operator of instance01 should not be able to view a cluster member",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Operator of instance01 should not be able to view a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Operator of instance01 should not be able to view a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of instance01 should not be able to view project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of instance01 should not be able to create instances in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_instances",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of instance01 should not be able to create images in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_images",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of instance01 should not be able to create networks in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of instance01 should not be able to create network ACLs in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_network_acls",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of instance01 should not be able to create network zones in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_network_zones",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of instance01 should not be able to create network forwards in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_network_forwards",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of instance01 should not be able to create network load balancers in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of instance01 should not be able to create network peers in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_network_peers",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of instance01 should not be able to create profiles in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_profiles",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of instance01 should not be able to create storage pool volumes in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_storage_pool_volumes",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of instance01 should not be able to create storage_buckets in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit an image in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "image:image01",
			},
		},
		{
			description: "Operator of instance01 should not be able to view an image in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "image:image01",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit instance01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of instance01 should be able to view instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of instance01 should be able to update the state of instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of instance01 should be able to manage snapshots of instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_manage_snapshots",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of instance01 should be able to manage backups of instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of instance01 should be able to connect to instance01 via sftp in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_connect_sftp",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of instance01 should be able to push/pull files into instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_access_files",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of instance01 should be able to access instance01's console",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_access_console",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of instance01 should be able to exec into instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_exec",
				Object:   "instance:instance01",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit a network in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "network:network01",
			},
		},
		{
			description: "Operator of instance01 should not be able to view a network in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "network:network01",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit a network ACL in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Operator of instance01 should not be able to view a network ACL in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit a network forward in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Operator of instance01 should not be able to view a network forward in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit a network load balancer in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Operator of instance01 should not be able to view a network load balancer in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit a network peer in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Operator of instance01 should not be able to view a network peer in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit a profile in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Operator of instance01 should not be able to view a profile in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "profile:profile01",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit a storage pool volume in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Operator of instance01 should not be able to view a storage pool volume in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "Operator of instance01 should not be able to edit a storage bucket in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_edit",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
		{
			description: "Operator of instance01 should not be able to view a storage bucket in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_operator",
				Relation: "can_view",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
	}

	for i, test := range tests {
		s.T().Logf("Case %d: %s", i, test.description)

		checkResponse, err := s.fga.Check(context.Background()).Options(client.ClientCheckOptions{AuthorizationModelId: s.authModelID}).Body(test.request).Execute()
		s.Require().NoError(err)
		s.Equal(test.allowed, checkResponse.GetAllowed())
	}
}

func (s *openFGASuite) TestInstanceUser() {
	clientWriteResponse, err := s.fga.WriteTuples(context.Background()).Options(client.ClientWriteOptions{AuthorizationModelId: s.authModelID}).Body(client.ClientWriteTuplesBody{
		{
			User:     "group:instance01_users#member",
			Relation: "user",
			Object:   "instance:instance01",
		},
		{
			User:     "user:instance01_user",
			Relation: "member",
			Object:   "group:instance01_users",
		},
	}).Execute()
	s.Require().NoError(err)

	s.Require().Len(clientWriteResponse.Deletes, 0)
	for _, write := range clientWriteResponse.Writes {
		s.Require().NoError(write.Error)
	}

	type test struct {
		description string
		allowed     bool
		request     client.ClientCheckRequest
	}

	tests := []test{
		{
			description: "User of instance01 should not be able to edit the server",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "User of instance01 should be able to view the server",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view_server",
				Object:   "server:lxd",
			},
		},
		{
			description: "User of instance01 should not be able to create a storage pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_storage_pool",
				Object:   "server:lxd",
			},
		},
		{
			description: "User of instance01 should not be able to create a project",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_project",
				Object:   "server:lxd",
			},
		},
		{
			description: "User of instance01 should not be able to view server resources",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view_resources",
				Object:   "server:lxd",
			},
		},
		{
			description: "User of instance01 should not be able to create a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_certificate",
				Object:   "server:lxd",
			},
		},
		{
			description: "User of instance01 should not be able to edit cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "User of instance01 should not be able to view cluster config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view_cluster",
				Object:   "server:lxd",
			},
		},
		{
			description: "User of instance01 should not be able to add a new member to a cluster",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_cluster_member",
				Object:   "server:lxd",
			},
		},
		{
			description: "User of instance01 should not be able to create a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_cluster_group",
				Object:   "server:lxd",
			},
		},
		{
			description: "User of instance01 should not be able to view server metrics",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view_metrics",
				Object:   "server:lxd",
			},
		},
		{
			description: "User of instance01 should not be able to edit a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "User of instance01 should not be able to view a certificate",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "certificate:eeef45f0570ce713864c86ec60c8d88f60b4844d3a8849b262c77cb18e88394d",
			},
		},
		{
			description: "User of instance01 should not be able to edit cluster member config",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "User of instance01 should not be able to view a cluster member",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "cluster_member:node01",
			},
		},
		{
			description: "User of instance01 should not be able to edit a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "User of instance01 should not be able to view a cluster group",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "cluster_group:group01",
			},
		},
		{
			description: "User of instance01 should not be able to edit a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "User of instance01 should not be able to view a storage_pool",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "storage_pool:pool01",
			},
		},
		{
			description: "User of instance01 should not be able to edit project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "project:project01",
			},
		},
		{
			description: "User of instance01 should not be able to view project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "project:project01",
			},
		},
		{
			description: "User of instance01 should not be able to create instances in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_instances",
				Object:   "project:project01",
			},
		},
		{
			description: "User of instance01 should not be able to create images in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_images",
				Object:   "project:project01",
			},
		},
		{
			description: "User of instance01 should not be able to create networks in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "User of instance01 should not be able to create network ACLs in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_network_acls",
				Object:   "project:project01",
			},
		},
		{
			description: "User of instance01 should not be able to create network zones in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_network_zones",
				Object:   "project:project01",
			},
		},
		{
			description: "User of instance01 should not be able to create network forwards in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_network_forwards",
				Object:   "project:project01",
			},
		},
		{
			description: "User of instance01 should not be able to create network load balancers in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "User of instance01 should not be able to create network peers in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_network_peers",
				Object:   "project:project01",
			},
		},
		{
			description: "User of instance01 should not be able to create profiles in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_profiles",
				Object:   "project:project01",
			},
		},
		{
			description: "User of instance01 should not be able to create storage pool volumes in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_storage_pool_volumes",
				Object:   "project:project01",
			},
		},
		{
			description: "User of instance01 should not be able to create storage_buckets in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_create_networks",
				Object:   "project:project01",
			},
		},
		{
			description: "User of instance01 should not be able to edit an image in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "image:image01",
			},
		},
		{
			description: "User of instance01 should not be able to view an image in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "image:image01",
			},
		},
		{
			description: "User of instance01 should not be able to edit instance01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User of instance01 should be able to view instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User of instance01 should not be able to update the state of instance01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User of instance01 should not be able to manage snapshots of instance01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_manage_snapshots",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User of instance01 should not be able to manage backups of instance01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_update_state",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User of instance01 should be able to connect to instance01 via sftp in project01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_connect_sftp",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User of instance01 should be able to push/pull files into instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_access_files",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User of instance01 should be able to access instance01's console",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_access_console",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User of instance01 should be able to exec into instance01",
			allowed:     true,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_exec",
				Object:   "instance:instance01",
			},
		},
		{
			description: "User of instance01 should not be able to edit a network in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "network:network01",
			},
		},
		{
			description: "User of instance01 should not be able to view a network in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "network:network01",
			},
		},
		{
			description: "User of instance01 should not be able to edit a network ACL in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "User of instance01 should not be able to view a network ACL in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "User of instance01 should not be able to edit a network forward in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "User of instance01 should not be able to view a network forward in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "network_acl:network_acl01",
			},
		},
		{
			description: "User of instance01 should not be able to edit a network load balancer in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "User of instance01 should not be able to view a network load balancer in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "network_load_balancer:network_load_balancer01",
			},
		},
		{
			description: "User of instance01 should not be able to edit a network peer in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "User of instance01 should not be able to view a network peer in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "network_peer:network_peer01",
			},
		},
		{
			description: "User of instance01 should not be able to edit a profile in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "profile:profile01",
			},
		},
		{
			description: "User of instance01 should not be able to view a profile in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "profile:profile01",
			},
		},
		{
			description: "User of instance01 should not be able to edit a storage pool volume in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "User of instance01 should not be able to view a storage pool volume in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "storage_pool_volume:storage_pool_volume01",
			},
		},
		{
			description: "User of instance01 should not be able to edit a storage bucket in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_edit",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
		{
			description: "User of instance01 should not be able to view a storage bucket in project01",
			allowed:     false,
			request: client.ClientCheckRequest{
				User:     "user:instance01_user",
				Relation: "can_view",
				Object:   "storage_bucket:storage_bucket01",
			},
		},
	}

	for i, test := range tests {
		s.T().Logf("Case %d: %s", i, test.description)

		checkResponse, err := s.fga.Check(context.Background()).Options(client.ClientCheckOptions{AuthorizationModelId: s.authModelID}).Body(test.request).Execute()
		s.Require().NoError(err)
		s.Equal(test.allowed, checkResponse.GetAllowed())
	}
}
