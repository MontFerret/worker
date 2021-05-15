#!/bin/bash
set -e

# Start mitm
exec sh -c "mitmdump -p 8080 -s '/inject.py'"