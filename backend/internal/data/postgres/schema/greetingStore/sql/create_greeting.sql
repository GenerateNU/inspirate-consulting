INSERT INTO public.greetings (
    name, greeting_string
) VALUES (
    $1, $2
)
RETURNING greeting_string