import Cookies from "js-cookie"
import React, { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import "./Login.css"

const LoginPage = () => {
	const [username, setUsername] = useState("")
	const [password, setPassword] = useState("")
	const [error, setError] = useState("")
	const navigate = useNavigate()

	useEffect(() => {
		var userRole = Cookies.get("X-User-Role")
		if (userRole === "admin") {
			navigate("/")
		}
	}, [])

	const handleLogin = e => {
		e.preventDefault()

		if (username === "admin" && password === "admin") {
			Cookies.set("X-User-Role", "admin", { path: "/", expires: 1 })
			navigate("/")
		} else {
			setError("Invalid username or password")
		}
	}

	return (
		<div className="login-container">
			<h1>Admin login</h1>
			<form onSubmit={handleLogin} className="login-form">
				<div className="form-group">
					<label htmlFor="username">Username:</label>
					<input
						type="text"
						id="username"
						value={username}
						onChange={e => setUsername(e.target.value)}
						required
					/>
				</div>
				<div className="form-group">
					<label htmlFor="password">Password:</label>
					<input
						type="password"
						id="password"
						value={password}
						onChange={e => setPassword(e.target.value)}
						required
					/>
				</div>
				{error && <p className="error-message">{error}</p>}
				<button type="submit" className="login-button">
					Login
				</button>
			</form>
		</div>
	)
}

export default LoginPage
