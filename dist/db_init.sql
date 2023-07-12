do $$ begin
   create type measuring_unit as enum (
      'gram',
      'mililitre',
      'unit'
   );
exception
   when duplicate_object then null;
end $$;

create table if not exists recipe (
   id int generated always as identity primary key,
   name varchar not null
);

create table if not exists ingredient (
   id int generated always as identity primary key,
   name varchar not null unique
);

create table if not exists recipe_ingredient (
   id int generated always as identity primary key,
   ingredient int not null,
   recipe int not null,
   amount numeric,
   measuring_unit measuring_unit not null,
   foreign key(ingredient) 
      references ingredient(id)
      on delete cascade,
   foreign key(recipe) 
      references recipe(id)
      on delete cascade
);
