
# Instructions to push to GitHub

1. Create a new repository on GitHub (e.g., `pgvms`).
2. Locally:
   git init
   git add .
   git commit -m "Initial PGVMS scaffold"
   git branch -M main
   git remote add origin git@github.com:YOUR_USERNAME/pgvms.git
   git push -u origin main

3. Add secrets and CI/CD as needed (cloud build service account, DB credentials).
