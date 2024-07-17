#!/bin/bash

template_file="docker-compose.template.yml"
output_file="docker-compose.yml"

# Copy the template to the output file
cp $template_file $output_file

# Generate node services
node_services=""
for i in {1..20}
do
  node_services=$(cat <<EOF
$node_services
  node$i:
    build:
      context: .
      args:
        P2P_CONFIG: node
    depends_on:
      - bootstrap
EOF
)
done

# Replace the placeholder with generated services
awk -v r="$node_services" '{gsub(/# NODE_SERVICES_PLACEHOLDER/, r)}1' $output_file > temp && mv temp $output_file

echo "docker-compose.yml has been generated."
