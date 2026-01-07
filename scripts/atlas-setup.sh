#!/bin/bash

# Atlas Setup Script
# This script helps you set up Atlas for the first time

set -e

echo "🚀 Atlas Migration Setup"
echo "======================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Atlas is installed
echo "1. Checking Atlas installation..."
if command -v atlas &> /dev/null; then
    echo -e "${GREEN}✓ Atlas is already installed${NC}"
    atlas version
else
    echo -e "${YELLOW}! Atlas not found. Installing...${NC}"
    curl -sSf https://atlasgo.sh | sh
    
    if command -v atlas &> /dev/null; then
        echo -e "${GREEN}✓ Atlas installed successfully!${NC}"
        atlas version
    else
        echo -e "${RED}✗ Failed to install Atlas${NC}"
        echo "Please install manually: https://atlasgo.io/getting-started#installation"
        exit 1
    fi
fi

echo ""
echo "2. Checking project structure..."

# Create migrations directory if not exists
if [ ! -d "migrations" ]; then
    echo "Creating migrations directory..."
    mkdir -p migrations
    echo -e "${GREEN}✓ Created migrations directory${NC}"
else
    echo -e "${GREEN}✓ migrations directory exists${NC}"
fi

# Check for atlas.hcl
if [ ! -f "atlas.hcl" ]; then
    echo -e "${RED}✗ atlas.hcl not found!${NC}"
    echo "Please ensure atlas.hcl exists in project root"
    exit 1
else
    echo -e "${GREEN}✓ atlas.hcl found${NC}"
fi

# Check for main.go in cmd/atlas
if [ ! -f "cmd/atlas/main.go" ]; then
    echo -e "${RED}✗ cmd/atlas/main.go not found!${NC}"
    echo "Please ensure cmd/atlas/main.go exists"
    exit 1
else
    echo -e "${GREEN}✓ cmd/atlas/main.go found${NC}"
fi

echo ""
echo "3. Checking environment configuration..."

# Check if .env exists
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}! .env file not found${NC}"
    echo "Please create .env file with database configuration"
    echo ""
    echo "Example for MySQL:"
    echo "ATLAS_DEV_DB_URL=mysql://user:pass@localhost:3306/dbname"
    echo ""
    echo "Example for PostgreSQL:"
    echo "ATLAS_DEV_DB_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable"
else
    echo -e "${GREEN}✓ .env file found${NC}"
    
    # Check if ATLAS_DEV_DB_URL is set
    if grep -q "ATLAS_DEV_DB_URL" .env; then
        echo -e "${GREEN}✓ ATLAS_DEV_DB_URL is configured${NC}"
    else
        echo -e "${YELLOW}! ATLAS_DEV_DB_URL not found in .env${NC}"
        echo "Add this to your .env file:"
        echo ""
        echo "# Atlas Migration Database URL"
        echo "ATLAS_DEV_DB_URL=mysql://user:pass@localhost:3306/dbname"
    fi
fi

echo ""
echo "4. Testing Atlas configuration..."

# Try to load schema
echo "Loading GORM schema..."
if go run ./cmd/atlas > /dev/null 2>&1; then
    echo -e "${GREEN}✓ GORM schema loads successfully${NC}"
else
    echo -e "${YELLOW}! Warning: Could not load GORM schema${NC}"
    echo "This might be normal if database connection is not configured yet"
fi

echo ""
echo "====================================="
echo -e "${GREEN}✅ Setup Complete!${NC}"
echo "====================================="
echo ""
echo "Next steps:"
echo ""
echo "1. Configure database URL in .env:"
echo "   ATLAS_DEV_DB_URL=mysql://user:pass@localhost:3306/dbname"
echo ""
echo "2. Generate your first migration:"
echo "   make atlas-diff"
echo ""
echo "3. Apply migrations:"
echo "   make atlas-apply"
echo ""
echo "4. Check migration status:"
echo "   make atlas-status"
echo ""
echo "📚 Documentation:"
echo "   - Project docs: docs/DATABASE.md"
echo "   - Migration README: migrations/README.md"
echo "   - Atlas docs: https://atlasgo.io/"
echo ""
echo "💡 Tip: Run 'make help' to see all available Atlas commands"
echo ""
