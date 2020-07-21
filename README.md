# Mailer
Simple, self-contained SMTP message sender

## Configuration

### environment variables

 - MAILER_HOST
 - MAILER_PORT
 - MAILER_USERNAME
 - MAILER_PASSWORD

### command line flags

 - --subject
 - --from
 - --to
 
## Usage

    echo "hello world" | mailer --to test@example.com --from test2@example.com --subject "hello"
