# OpenFGA in LXD

## Usage
Run the [OpenFGA server with docker](https://openfga.dev/docs/getting-started/setup-openfga/docker):
```shell
docker run -p 8080:8080 -p 8081:8081 -p 3000:3000 openfga/openfga run
```

Run the tests:
```shell
go test -v .
```

To iterate, edit the model in `lxd.openfga`, then run `make update-openfga`, and re-run the tests.

## Existing model proposal
Specification: https://discuss.linuxcontainers.org/t/lxd-rebac-authorization-using-openfga/17094#authorization-model-5
1. Follows current RBAC model closely (except for adding more fine-grained permissions for network ACLs and network zones).
2. Relies on the creation of an "admin" group which is not stated in model.
3. Following the guides on the [OpenFGA website](https://openfga.dev/docs/modeling/getting-started#introduction-to-modeling),
in most cases we should be able to check permissions with a single question: *Can user **U** perform action **A** on object **O**?*
However, in the current model we need to check: *Does user **U** have relation **member** with group **admin OR** can user **U** perform action **A** on object **O**?*
4. Juju model has 157 entitlements (relations). LXD may not be quite as complicated as Juju but the current model is not a fine-grained as it maybe should be.

## Proposed model
See [`lxd.openfga`](./lxd.openfga). In the Canonical terminology:
* Resources: These are `types`.
* Entitlements: Any relation that looks like `can_do_thing`.
* Roles: These are relations like `admin`, `manager`, or `operator` that are referenced by other relations and relations of child resources.
* Users: This is also a `type` that represents a single user. Users can be granted direct access on all entitlements.
* Groups: This is a `type` that has a direct relation `member` to `user`. Groups can also be granted direct access on all entitlements.

Some key points:
* A top level type `server` is created, representing a LXD server or cluster.
* The relation `admin` defined on `server` grants edit permissions on every resource in the cluster.
* The relation `operator` defined on `server` grants permissions to create and manage projects, but not to edit the cluster itself.
* The relation `viewer` defined on `server` grants permission to view all cluster resources.
* The relation `user` is a [type bound public access](https://openfga.dev/docs/concepts#what-is-type-bound-public-access) for any authenticated user.
This should be used for the `/1.0` endpoint (`server:can_view_server`).
* There are then a number of relations allowing creation of specific resources. For example, a group with `server:operator` privileges could be granted access to manage certificates with the tuple:
```
User: group:operators#member
Relation: can_create_certificates
Object: server:lxd
```
* Generally, if a resource does not have any child resources then it has `manager`, `viewer`, `can_edit`, and `can_view`. `manager` and `viewer` are added to allow 
access to a resource via a [direct access](https://openfga.dev/docs/modeling/direct-access). `can_edit` and `can_view` are the relations we will perform checks against.
These relations usually contain an inherited relation from a parent.
* The `project` type follows a similar pattern to `server` in that there is a `manager` who can edit the project and all resources within, 
an `operator` who can create resources within the project but not edit the project configuration, and a `viewer` who can view all project resources but not edit them.
There are also a number of relations for creating specific kinds of resources within a project.
* When permitting a user or group the `manager` or `operator` relation on a project, we also need to give them `server:viewer` so that they can view say a cluster group and be able to target it when creating an instance.
* The `instance` type contains a number of extra relations representing actions that can be performed on an instance:
  1. `manager` can perform any action on an instance.
  2. `operator` can change the instance state and manage backups/snapshots, but cannot edit instance config.
  3. `viewer` can view the instance config.
  4. `user` can interact with the instance via file push/pull, sftp, console, and exec. (E.g. ssh access but better).
  
### Use cases
* `server:admin` creates a project and grants a group `project:operator` permission on that project (plus `server:viewer`). 
Members of the group can create and manage resources in the project but cannot change project configuration (which could escalate privileges).
* `project:operator` creates an instance and grants a user `instance:user` permission. The user can connect to the instance but not edit it.

## General questions
1. What happens when the authentication method changes?
2. If not using Canonical OpenFGA, how does an administrator change permissions for a user or group?

## Questions about proposed model
1. What name do we give the top-level `server` object? Or is there a way to make it singular?
2. What to do about operations and warnings?
   a. Which users can cancel an operation? (My guess is `project:operator` but maybe we just have it as `server:user`, since you must have received a UUID from a protected endpoint?
   b. Should all users be able to view operations?
   c. Who needs to view cluster warnings? Who needs to access project specific warnings?