alter table public.media
    alter column school_year type year
    using case school_year
        when 9 then 'freshman'
        when 10 then 'sophomore'
        when 11 then 'junior'
        when 12 then 'senior'
    end::year;