import React from "react"
import Cookies from "js-cookie"
import { Link, useNavigate } from "react-router-dom"
import "./Navbar.css" // Импорт файла CSS

const Navbar = () => {
	const navigate = useNavigate()
	const userRole = Cookies.get("X-User-Role")

	const handleLogout = () => {
		Cookies.remove("X-User-Role")
		navigate("/login")
	}

	return (
		<nav className="navbar">
			<h1 className="navbar-title">Weather App</h1>
			<div className="navbar-links">
				<Link to="/" className="navbar-link">
					View Weather
				</Link>
				{userRole === "admin" && (
					<Link to="/add" className="navbar-link">
						Add Weather
					</Link>
				)}
				{userRole ? (
					<button onClick={handleLogout} className="navbar-button">
						Logout
					</button>
				) : (
					<Link to="/login" className="navbar-link">
						Login
					</Link>
				)}
			</div>
		</nav>
	)
}

export default Navbar
