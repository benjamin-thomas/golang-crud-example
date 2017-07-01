#!/usr/bin/env ruby

require 'bundler'
require 'logger'
Bundler.require

user = ENV.fetch('PGUSER')
host = `docker inspect --format '{{ .NetworkSettings.IPAddress  }}' golang-crud-example-pg`.chomp
pw = ENV.fetch('PGPASSWORD')
port = 5432
db_name = ENV.fetch('PGDATABASE')

DB = Sequel.connect("postgres://#{user}:#{pw}@#{host}:#{port}/#{db_name}")
DB.loggers << Logger.new(STDOUT)

countries_count = DB['SELECT COUNT(*) FROM countries'].first.fetch(:count)

countries = DB.from(:countries)

puts "Country count : #{countries_count}"
puts '---'
puts "Country columns: #{countries.columns}"
puts '---'

puts 'Random countries'
puts '================'
countries.order(Sequel.function(:random)).limit(10).each do |c|
  printf "%03d: #{c.fetch(:name)}\n", c.fetch(:id)
end

puts countries.where(name: 'France').limit(1).first.fetch(:id)
puts countries[name: 'France'].fetch(:id)

res = countries.where(name: 'France').limit(1).select(:id)
puts res.sql

puts countries.where(Sequel.ilike(:name, 'france')).first.fetch(:id)
puts countries[Sequel.ilike(:name, 'france')].fetch(:id)
puts countries.first(Sequel.ilike(:name, 'france')).fetch(:id)

p countries.limit(3).to_hash(:id, :name)
p countries.limit(10).to_hash_groups(:created_on, :id)
require 'pry' ; binding.pry
