# UP Namespace Security & Versioning

Namespace plugins require strong security guarantees and reproducible builds. UP provides cryptographic verification, version pinning, and compatibility checking.

## Version Specification

### Explicit Version Pinning

```up
# Exact version required
!use [time@1.2.3, random@2.0.0, faker@3.1.4]

# Version range (semantic versioning)
!use [time>=1.0.0,<2.0.0]

# Compatible versions (^)
!use [greeting^1.2.0]  # Allows 1.2.0 to <2.0.0

# Minimum version (~)
!use [greeting~1.2.0]  # Allows 1.2.x

# Latest compatible
!use [greeting@latest]
```

### Version Pinning with Aliases

When using aliases, versions are specified before the alias:

```up
# Multiple versions of the same namespace
!use [
  time@1.0.0 as oldtime,
  time@2.0.0 as newtime
]

created_at $oldtime.now      # Uses time 1.0.0
updated_at $newtime.now      # Uses time 2.0.0
```

```up
# Version requirements with different sources
!use [
  github.com/uplang/ns-time@1.5.0 as time,
  github.com/myorg/time@2.0.0 as customtime
]

# Both are pinned to specific versions
standard $time.now
custom $customtime.now
```

```up
# Mixing default and aliased
!use [
  time,                          # Default latest
  time@1.0.0 as legacytime,     # Pin old version
  time^2.0.0 as moderntime      # Compatible with 2.x
]
```

### Version Requirements in Schema

**File: `greeting.up-schema`**

```up
namespace greeting

version 1.2.3

# Compatibility information
compatibility {
  breaking_changes [2.0.0]
  deprecated_since 1.5.0
  min_up_version 1.0.0
}

# Dependencies on other namespaces
depends {
  time >= 1.0.0
  id ^2.0.0
}
```

## Security Model

### Hash Verification

Every namespace executable must have a verified hash. UP uses SHA-256 for integrity checking.

**File: `up-namespaces.lock`** (automatically generated)

```up
# UP Namespace Lock File
# This file pins exact versions and hashes for reproducibility
# DO NOT EDIT MANUALLY - managed by up namespace commands

version 1.0.0
generated_at 2025-10-05T14:30:00Z

namespaces {
  greeting {
    version 1.2.3

    files {
      executable {
        path ./up-namespaces/greeting
        hash sha256:a3b5c9d1e2f4567890abcdef1234567890abcdef1234567890abcdef12345678
        size!int 1195
      }

      schema {
        path ./up-namespaces/greeting.up-schema
        hash sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
        size!int 2286
      }
    }

    verified!bool true
    verified_at 2025-10-05T14:30:00Z
    source local
  }

  time {
    version 2.1.0

    builtin!bool true
    verified!bool true
    source builtin
  }

  faker {
    version 3.1.4

    files {
      executable {
        path ~/.up/namespaces/faker
        hash sha256:9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba
        size!int 45120
      }

      schema {
        path ~/.up/namespaces/faker.up-schema
        hash sha256:fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210
        size!int 8945
      }
    }

    verified!bool true
    verified_at 2025-10-05T14:29:00Z
    source registry
    registry_url https://registry.uplang.org/faker/3.1.4
  }
}

# Aliases mapping
aliases {
  oldtime time@1.0.0
  newtime time@2.1.0
  customtime github.com/myorg/time@2.0.0
}
```

### Lock File with Aliases

When using namespace aliases, the lock file tracks both the alias and the actual namespace:

**File: `up-namespaces.lock`**

```up
version 1.0.0
generated_at 2025-10-05T15:00:00Z

# Template uses: !use [time@1.0.0 as oldtime, time@2.0.0 as newtime]

namespaces {
  time {
    version 1.0.0
    alias oldtime

    files {
      executable {
        path ~/.up/namespaces/time/1.0.0/time
        hash sha256:abc123...
        size!int 12400
      }
    }

    verified!bool true
    source registry
  }

  time {
    version 2.0.0
    alias newtime

    files {
      executable {
        path ~/.up/namespaces/time/2.0.0/time
        hash sha256:def456...
        size!int 15200
      }
    }

    verified!bool true
    source registry
  }

  customtime {
    namespace github.com/myorg/time
    version 2.0.0

    files {
      executable {
        path ~/.up/namespaces/github.com/myorg/time/2.0.0/time
        hash sha256:789xyz...
        size!int 18900
      }
    }

    verified!bool true
    source url
    url https://github.com/myorg/time
  }
}

# Alias resolution table
aliases {
  oldtime time@1.0.0
  newtime time@2.0.0
  customtime github.com/myorg/time@2.0.0
}
```

### Cryptographic Signatures

Namespaces can be cryptographically signed for trust and authenticity.

**File: `greeting.up-signature`**

```up
# UP Namespace Signature
# Cryptographic signature for namespace verification

namespace greeting
version 1.2.3

signature {
  algorithm ed25519
  public_key 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
  signature 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
  signed_at 2025-10-05T14:00:00Z
}

files {
  executable {
    path greeting
    hash sha256:a3b5c9d1e2f4567890abcdef1234567890abcdef1234567890abcdef12345678
  }

  schema {
    path greeting.up-schema
    hash sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
  }
}

signer {
  name UP Examples Team
  email security@up-examples.org
  keybase up_examples
  github nicerobot
}
```

## Security Verification Process

### 1. Version Resolution

```
User specifies: !use [greeting^1.2.0]
↓
UP resolves: greeting@1.2.3 (latest compatible)
↓
Checks up-namespaces.lock for pinned version
↓
If no lock: downloads/locates greeting@1.2.3
If locked: uses exact version from lock file
```

### 2. Hash Verification

```
1. Read executable from disk
2. Calculate SHA-256 hash
3. Compare with hash in up-namespaces.lock
4. If mismatch: FAIL with security error
5. If match: proceed to signature verification
```

### 3. Signature Verification (Optional but Recommended)

```
1. Check for .up-signature file
2. Load public key from signature
3. Verify signature against file hashes
4. Check signer identity against trusted keys
5. If invalid: FAIL or WARN (based on policy)
6. If valid: mark as verified
```

### 4. Compatibility Check

```
1. Load namespace schema
2. Check version compatibility rules
3. Verify dependencies
4. Check UP version requirements
5. If incompatible: FAIL with clear error
6. If compatible: allow execution
```

## Security Policy File

**File: `.up-security`** (project-level security policy)

```up
# UP Security Policy
# Defines security requirements for namespace plugins

version 1.0.0

policy {
  # Require signatures for non-local namespaces
  require_signatures!bool true

  # Allow only these signature algorithms
  allowed_algorithms [ed25519, rsa4096]

  # Verify hashes on every execution
  verify_hashes!bool true

  # Trust level required
  min_trust_level verified  # none, local, verified, official

  # Allow unsigned local namespaces
  allow_unsigned_local!bool true

  # Fail on hash mismatch
  fail_on_hash_mismatch!bool true

  # Warn on deprecated versions
  warn_on_deprecated!bool true
}

# Trusted signers (public keys)
trusted_signers {
  up_core {
    public_key 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
    name UP Core Team
    trust_level official
  }

  community_verified {
    public_key 0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321
    name UP Community
    trust_level verified
  }
}

# Namespace-specific policies
namespace_policies {
  # Built-in namespaces are always trusted
  builtin {
    namespaces [time, date, id, random]
    trust_level official
    require_signature!bool false
  }

  # Dangerous namespaces require explicit approval
  dangerous {
    namespaces [file, env, exec]
    require_explicit_approval!bool true
    require_signature!bool true
    min_trust_level official
  }

  # Local development namespaces
  local {
    namespaces_pattern ^local_.*
    require_signature!bool false
    trust_level local
  }
}

# Audit logging
audit {
  enabled!bool true
  log_path ~/.up/audit.log

  log_events [
    namespace_load
    signature_verification
    hash_mismatch
    version_incompatibility
    security_violation
  ]
}
```

## CLI Commands for Security

### Initialize Security

```bash
# Create up-namespaces.lock
up namespace lock

# Generate lock file with current namespace versions
up namespace lock --update

# Verify all namespaces match lock file
up namespace verify
```

### Signature Management

```bash
# Generate signing key pair
up namespace keygen --name "My Team" --email "security@example.com"

# Sign a namespace
up namespace sign greeting --key ~/.up/keys/private.key

# Verify signature
up namespace verify-signature greeting

# Add trusted signer
up namespace trust add 0x1234...cdef --name "UP Core"

# List trusted signers
up namespace trust list

# Remove trusted signer
up namespace trust remove 0x1234...cdef
```

### Version Management

```bash
# Show installed versions
up namespace list --versions

# Update to latest compatible versions
up namespace update

# Update specific namespace
up namespace update greeting

# Pin to specific version
up namespace pin greeting@1.2.3

# Show outdated namespaces
up namespace outdated
```

### Security Auditing

```bash
# Audit all namespaces
up namespace audit

# Check for security issues
up namespace security-check

# View audit log
up namespace audit-log

# Generate security report
up namespace security-report -o report.up
```

## Example Security Workflow

### 1. Project Setup

```bash
# Initialize new project
cd myproject

# Create security policy
cat > .up-security << 'EOF'
policy {
  require_signatures!bool true
  verify_hashes!bool true
}
EOF

# Install namespaces with verification
up namespace install greeting@1.2.3 --verify
up namespace install faker@3.1.4 --verify
```

### 2. Lock File Generated

**`up-namespaces.lock`** created automatically:

```up
version 1.0.0
generated_at 2025-10-05T14:30:00Z

namespaces {
  greeting {
    version 1.2.3
    files {
      executable {
        hash sha256:a3b5c9d1e2f4567890abcdef1234567890abcdef1234567890abcdef12345678
      }
    }
    verified!bool true
  }

  faker {
    version 3.1.4
    files {
      executable {
        hash sha256:9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba
      }
    }
    verified!bool true
  }
}
```

### 3. Using in UP

```up
!use [greeting@1.2.3, faker@3.1.4]

# UP verifies:
# 1. Exact versions match lock file
# 2. File hashes match lock file
# 3. Signatures are valid (if required)
# 4. Versions are compatible

test_user {
  name $faker.name
  greeting $greeting.hello($faker.name)
}
```

### 4. Verification on Execution

```
$ up template process -i test.up

[INFO] Loading namespaces...
[INFO] Verifying greeting@1.2.3...
  ✓ Version matches lock file
  ✓ Hash verified: sha256:a3b5c9d1e2...
  ✓ Signature valid (UP Examples Team)
[INFO] Verifying faker@3.1.4...
  ✓ Version matches lock file
  ✓ Hash verified: sha256:9876543210...
  ✓ Signature valid (UP Core Team)
[SUCCESS] All namespaces verified

[OUTPUT]
test_user {
  name Alice Johnson
  greeting Hello, Alice Johnson!
}
```

### 5. Security Failure Example

```
$ up template process -i test.up

[INFO] Loading namespaces...
[INFO] Verifying greeting@1.2.3...
  ✗ Hash mismatch!
    Expected: sha256:a3b5c9d1e2f4567890abcdef...
    Got:      sha256:deadbeef1234567890abcdef...

[ERROR] Security verification failed for namespace: greeting
[ERROR] File has been modified or corrupted
[ERROR] Refusing to execute

Suggestion:
  1. Re-install the namespace: up namespace install greeting@1.2.3 --force
  2. Update lock file: up namespace lock --update
  3. Report security issue if unexpected
```

## Registry Integration

Namespaces can be published to a central registry with automatic signing.

### Publishing to Registry

```bash
# Package namespace for distribution
up namespace package greeting

# Publish to registry (requires authentication)
up namespace publish greeting --registry https://registry.uplang.org

# Package includes:
# - executable
# - schema
# - signature
# - README
# - LICENSE
```

### Registry Manifest

**Registry stores:**

```up
# registry.uplang.org/greeting/manifest.up

namespace greeting

versions {
  1.2.3 {
    released_at 2025-10-05T14:00:00Z

    files {
      executable {
        url https://registry.uplang.org/greeting/1.2.3/greeting
        hash sha256:a3b5c9d1e2f4567890abcdef1234567890abcdef1234567890abcdef12345678
        size!int 1195
      }

      schema {
        url https://registry.uplang.org/greeting/1.2.3/greeting.up-schema
        hash sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
        size!int 2286
      }

      signature {
        url https://registry.uplang.org/greeting/1.2.3/greeting.up-signature
        algorithm ed25519
        public_key 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
      }
    }

    verified!bool true
    verified_by UP Registry
    downloads!int 12845

    dependencies {
      up >= 1.0.0
    }

    checksums {
      sha256 a3b5c9d1e2f4567890abcdef1234567890abcdef1234567890abcdef12345678
      sha512 fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210
    }
  }

  1.2.2 {
    released_at 2025-09-15T10:00:00Z
    deprecated!bool true
    deprecated_reason Security vulnerability fixed in 1.2.3
    # ... similar structure
  }
}

metadata {
  author UP Examples Team
  repository https://github.com/up-examples/greeting
  license MIT
  homepage https://up-examples.org/greeting

  tags [greeting, example, simple]

  stats {
    total_downloads!int 45678
    monthly_downloads!int 3456
    stars!int 234
  }
}
```

### Installing from Registry

```bash
# Install from registry with automatic verification
up namespace install greeting@1.2.3

# Process:
# 1. Query registry for greeting@1.2.3
# 2. Download files (executable, schema, signature)
# 3. Verify checksums against manifest
# 4. Verify signature with public key
# 5. Install to ~/.up/namespaces/
# 6. Update up-namespaces.lock
# 7. Mark as verified
```

## Security Best Practices

### For Namespace Developers

1. **Always sign releases** - Use `up namespace sign`
2. **Pin dependencies** - Specify exact versions in schema
3. **Document security** - Note any special requirements
4. **Version semantically** - Follow semver for compatibility
5. **Test thoroughly** - Security bugs affect all users
6. **Minimize permissions** - Request only needed capabilities
7. **Audit regularly** - Review for vulnerabilities

### For Namespace Users

1. **Pin versions** - Use exact versions in production
2. **Verify signatures** - Enable signature checking
3. **Review before use** - Check namespace code and schema
4. **Use lock files** - Commit `up-namespaces.lock`
5. **Update carefully** - Test updates before deploying
6. **Monitor audit logs** - Watch for security events
7. **Report issues** - Security bugs should be reported

### For Organizations

1. **Private registry** - Host internal namespaces
2. **Require signatures** - Enforce signing policy
3. **Audit regularly** - Review namespace usage
4. **Whitelist namespaces** - Approve before use
5. **Automated scanning** - Check for vulnerabilities
6. **Security training** - Educate developers
7. **Incident response** - Plan for security issues

## Threat Model

### Threats Mitigated

✅ **Supply chain attacks** - Cryptographic verification prevents tampering
✅ **Version confusion** - Lock files ensure reproducibility
✅ **Malicious updates** - Signature verification blocks unauthorized versions
✅ **Code injection** - Hash verification detects modifications
✅ **Dependency confusion** - Registry verification ensures authenticity

### Threats Requiring Additional Controls

⚠️ **Compromised signing keys** - Key rotation and revocation needed
⚠️ **Zero-day in namespace** - Regular security audits required
⚠️ **Malicious registry** - Use trusted registries only
⚠️ **Privilege escalation** - Sandboxing and capability limits
⚠️ **Side-channel attacks** - Secure execution environment

## Summary

**Version Specification:**
- ✅ Semantic versioning with range operators
- ✅ Exact pinning via lock files
- ✅ Compatibility declarations in schemas

**Security:**
- ✅ SHA-256 hash verification
- ✅ Ed25519/RSA cryptographic signatures
- ✅ Trusted signer management
- ✅ Security policy files

**Verification Process:**
- ✅ Version resolution
- ✅ Hash checking
- ✅ Signature verification
- ✅ Compatibility validation

**Tools:**
- ✅ `up namespace lock` - Pin versions
- ✅ `up namespace verify` - Check integrity
- ✅ `up namespace sign` - Cryptographic signing
- ✅ `up namespace audit` - Security review

**UP ensures namespace plugins are secure, reproducible, and trustworthy.**

