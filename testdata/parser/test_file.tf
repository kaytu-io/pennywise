locals {
  proxy = var.cats_mother
}

variable "cats_mother" {
  default = "boots"
}

provider "cats" {

}

resource "cats_cat" "mittens" {
  name = "mittens"
  special = true
}

resource "cats_kitten" "the-great-destroyer" {
  name = "the great destroyer"
  parent = cats_cat.mittens.name
}

data "cats_cat" "the-cats-mother" {
  name = local.proxy
}