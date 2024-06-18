# Deployinator

Ever used a IaaS provider like Vercel, Netlify or Railway and loved how automagically your code gets deployed as soon as you push to the repository? Well, now you can do the same with **Deployinator**, a lightweight deployment server built using the Gin web framework for Go.

It listens for incoming deployment requests from GitHub Webhooks, validates them using HMAC signatures, and triggers corresponding deployment scripts for different projects.

## Features

- Secure deployment using HMAC SHA256 signatures
- Easy-to-extend with deployment scripts
- You'll never have to worry about a <a href="https://serverlesshorrors.com/" target="_blank">random $10k bill ever again</a>.

## Installation

### Prerequisites

- Go (version 1.21 or above)
- Bash

### Steps

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/deployinator.git
   cd deployinator
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Create a `.env` file with the following content:

   ```env
   GIN_MODE=release
   SECRET_KEY=super_secret_key
   ```

4. Create your deployment scripts inside the `deploy_scripts` directory. Each script should be named after the project it deploys, e.g., `my_project.sh`.

5. Build and run the application:

   ```sh
   go build -o deployinator
   ./deployinator
   ```

## Configuration

### Environment Variables

- `GIN_MODE`: Sets the Gin mode. Options are `debug`, `release`, and `test`.
- `SECRET_KEY`: A secret key used to validate incoming deployment requests from GitHub.

## Usage

1. Deployinator listens on port `4444` by default.
2. To trigger a deployment, GitHub Webhooks should be configured to send a POST request to `http://yourdomain.com/<projectName>` with the HMAC SHA256 signature header.

### Setting up GitHub Webhooks

1. Go to your GitHub repository and navigate to **Settings** > **Webhooks**.
2. Click **Add webhook**.
3. In the **Payload URL** field, enter `http://yourdomain.com/<projectName>`, replacing `<projectName>` with the name of your project.
4. Set the **Content type** to `application/x-www-form-urlencoded`.
5. In the **Secret** field, enter the `SECRET_KEY` from your `.env` file.
6. Choose the events you want to trigger the webhook. Typically, **Push events** are used for deployments.
7. Click **Add webhook** to save your settings.

### Deployment Scripts

Deployment scripts should be placed in the `deploy_scripts` directory and should be named according to the project they deploy. For example, a script for the project `my_project` should be named `my_project.sh`.

#### Example Script

```sh
#!/bin/bash

# Example deployment script for my_project

echo "Deploying my_project"
# Add your deployment commands here
```

Make sure your scripts are executable:

```sh
chmod +x deploy_scripts/*.sh
```

## Security

Deployinator uses HMAC SHA256 to verify the authenticity of the deployment requests from GitHub. Ensure that your `SECRET_KEY` is kept secure and never exposed in your code or logs.

## License

Deployinator is open-sourced software licensed under the [MIT license](LICENSE).

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss your changes.
