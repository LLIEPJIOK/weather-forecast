import React from "react"
import { Link } from "react-router-dom"
import "./Navbar.css" // Импорт файла CSS

const Navbar = () => {
	return (
		<nav className="navbar">
			<h1 className="navbar-title">Weather App</h1>
			<div className="navbar-links">
				<Link to="/" className="navbar-link">
					View Weather
				</Link>
				<Link to="/add" className="navbar-link">
					Add Weather
				</Link>
			</div>
		</nav>
	)
}

export default Navbar
