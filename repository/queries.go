package repository

const (
	InsertUserQuery = `
		INSERT INTO users (phone, name, password) 
		VALUES ($1, $2, $3)
		RETURNING id`

	GetUserByIDQuery = `
		SELECT
			id,
			phone,
			name,
			password,
			created_at,
			updated_at
		FROM
			users
		WHERE id = $1`

	GetUserByPhoneQuery = `
		SELECT
			id,
			phone,
			name,
			password,
			created_at,
			updated_at
		FROM
			users
		WHERE phone = $1`

	UpdateUserQuery = `
		UPDATE users
		SET
			phone = CASE WHEN $1 != '' THEN $1 ELSE phone END,
			name = CASE WHEN $2 != '' THEN $2 ELSE name END,
			updated_at = now()
		WHERE id = $3`

	UpdateLoginCount = `
		UPDATE users
		SET
			login_count = login_count + 1
		WHERE id = $1`
)
