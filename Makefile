.PHONY: backend frontend run push clean

backend:
	cd backend && go run main.go

frontend:
	cd frontend && \
	if [ ! -d "node_modules" ]; then npm install; fi && \
	npm run dev

run:
	(cd backend && go run main.go) & \
	(cd frontend && \
		if [ ! -d "node_modules" ]; then npm install; fi && \
		npm run dev)

push:
	git add . && \
	read -r -p "Your commit message: " message && \
	git commit -m "$$message" && \
	git push

clean:
	rm -f backend/internal/database/*db
	rm -rf backend/media
