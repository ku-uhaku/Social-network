.PHONY: backend frontend install run

backend:
	cd backend && go run main.go

frontend:
	cd frontend && \
	if [ ! -d "node_modules" ]; then npm install; fi && \
	npm run dev

install:
	cd frontend && npm install


run:
	(cd backend && go run main.go) & \
	(cd frontend && \
		if [ ! -d "node_modules" ]; then npm install; fi && \
		npm run dev)
