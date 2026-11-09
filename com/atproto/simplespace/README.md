# com.atproto.simplespace

An implementation of the (deprecated) [`com.atproto.simplespace`](https://github.com/bluesky-social/proposals/blob/a0e3a5be505dd9d28f689026a882dd886b356a90/0016-permissioned-data/README.md#required-pds-space-management-simplespace) XRPC API defined here:

| Method                                 | Type      | Description                                                                   |
|----------------------------------------|-----------|-------------------------------------------------------------------------------|
| `com.atproto.simplespace.createSpace`  | procedure | Create an ATProto space. (Caller becomes the authority)                       |
| `com.atproto.simplespace.updateSpace`  | procedure | Update config.                                                                |
| `com.atproto.simplespace.deleteSpace`  | procedure | Delete the ATProto space.                                                     |
| `com.atproto.simplespace.getSpace`     | query     | Describe an ATProto space and its configuration.                              |
| `com.atproto.simplespace.putMember`    | procedure | Add a member (by DID) or update their read and write access.                  |
| `com.atproto.simplespace.removeMember` | procedure | Remove a member (by DID).                                                     |
| `com.atproto.simplespace.listMembers`  | query     | List the current members of an ATProto space and their read and write access. |
