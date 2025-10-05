# UP Namespace: fake

Fake data generation for testing, prototyping, and development.

## Installation

```bash
go install github.com/uplang/ns-fake@latest
```

Or build from source:

```bash
cd ns/fake
go mod download
go build -o fake
```

## Usage in UP Templates

```up
!use [fake]

# Generate test users
users $list.generate(count=10, template={
  id $id.uuid
  name $fake.name
  email $fake.email
  phone $fake.phone
  username $fake.username
})

# Company data
company {
  name $fake.company
  address $fake.address
  city $fake.city
  state $fake.state
  country $fake.country
  zipCode $fake.zipCode
}

# Product catalog
products $list.generate(count=50, template={
  name $fake.product
  price!float $fake.price(min=9.99, max=999.99)
  currency $fake.currency
  color $fake.color
})

# Web data
website {
  url $fake.url
  domain $fake.domain
  ipv4 $fake.ipv4
  userAgent $fake.userAgent
}

# Text content
content {
  headline $fake.sentence(words=5)
  description $fake.paragraph(sentences=3)
  body $fake.lorem(words=200)
}
```

## Functions

### Person Functions

- **`name()`** - Full name
- **`firstName()`** - First name only
- **`lastName()`** - Last name only
- **`email()`** - Email address
- **`phone()`** - Phone number
- **`username()`** - Username

### Internet Functions

- **`url()`** - Full URL
- **`domain()`** - Domain name
- **`ipv4()`** - IPv4 address
- **`ipv6()`** - IPv6 address
- **`userAgent()`** - Browser user agent string

### Company Functions

- **`company()`** - Company name
- **`jobTitle()`** - Job title/position

### Address Functions

- **`address()`** - Street address
- **`city()`** - City name
- **`state()`** - State/province name
- **`country()`** - Country name
- **`zipCode()`** - ZIP/postal code
- **`latitude()`** - Latitude coordinate
- **`longitude()`** - Longitude coordinate

### Text Functions

- **`word()`** - Single word
- **`sentence(words?)`** - Sentence (default: 10 words)
- **`paragraph(sentences?)`** - Paragraph (default: 3 sentences)
- **`lorem(words?)`** - Lorem ipsum text (default: 50 words)

### Commerce Functions

- **`product()`** - Product name
- **`price(min?, max?)`** - Price (default: $1.00-$1000.00)
- **`currency()`** - Currency code (USD, EUR, etc.)

### Color Functions

- **`color()`** - Color name
- **`hexColor()`** - Hex color code (#RRGGBB)

### Misc Functions

- **`creditCard(type?)`** - Credit card number (for testing only!)

## Use Cases

### Database Seeding

```up
!use [fake, id]

users $list.generate(count=1000, template={
  id $id.uuid
  first_name $fake.firstName
  last_name $fake.lastName
  email $fake.email
  phone $fake.phone
  company $fake.company
  job_title $fake.jobTitle
  address $fake.address
  city $fake.city
  state $fake.state
  zip_code $fake.zipCode
  created_at $time.now
})
```

### API Mocking

```up
!use [fake, id, random]

api_response {
  status 200
  data {
    user {
      id $id.uuid
      name $fake.name
      email $fake.email
      avatar $fake.url
    }
    session {
      token $random.bytes(size=32)
      expires_at $time.add(duration="24h")
    }
  }
}
```

### UI Prototyping

```up
!use [fake]

blog_posts $list.generate(count=20, template={
  title $fake.sentence(words=6)
  excerpt $fake.paragraph(sentences=2)
  body $fake.lorem(words=500)
  author {
    name $fake.name
    avatar $fake.url
  }
  tags $list.generate(count=3, template=$fake.word)
})
```

### Test Data

```up
!use [fake, random]

test_scenarios [
  {
    name "valid user registration"
    input {
      name $fake.name
      email $fake.email
      password $random.bytes(size=16)
      age!int $random.int(min=18, max=99)
    }
    expected_status 201
  }
  {
    name "international user"
    input {
      name $fake.name
      country $fake.country
      phone $fake.phone
      language $fake.word
    }
    expected_status 201
  }
]
```

## Seeding for Reproducibility

While not yet implemented, seeding will allow consistent data generation:

```up
# Future: Seeded generation for reproducible tests
!use [fake]
!seed 12345

user {
  name $fake.name  # Always generates same name with same seed
  email $fake.email
}
```

## Security Note

**Warning:** The `creditCard()` function generates fake credit card numbers for **testing purposes only**. Never use these for real financial transactions or store them as if they were real.

## Testing

```bash
echo '{"function":"name","params":{},"context":{}}' | ./fake
echo '{"function":"email","params":{},"context":{}}' | ./fake
echo '{"function":"address","params":{},"context":{}}' | ./fake
```

## Dependencies

This namespace uses [jaswdr/faker](https://github.com/jaswdr/faker) for data generation.

## License

MIT License

## Links

- [UP Language](https://uplang.org)
- [UP Specification](https://github.com/uplang/spec)
- [Schema](fake.up-schema)
