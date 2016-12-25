$module('generic', function () {
  var $static = this;
  this.$class('989e7835', 'SomeClass', true, '', function (T) {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.Something = function () {
      var $this = this;
      return null;
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
      };
      computed[("Something|2|29dc432d<" + $t.typeid(T)) + ">"] = true;
      return this.$cachedtypesig = computed;
    };
  });

  this.$class('e0aeb390', 'A', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  this.$class('2989c0d0', 'B', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  this.$interface('9be13b13', 'ASomething', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Something|2|29dc432d<e0aeb390>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$interface('7c58da5e', 'BSomething', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Something|2|29dc432d<2989c0d0>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var asc;
    var asc2;
    var bsc;
    asc = $g.generic.SomeClass($g.generic.A).new();
    asc2 = $g.generic.SomeClass($g.generic.A).new();
    bsc = $g.generic.SomeClass($g.generic.B).new();
    $t.cast(asc, $g.generic.ASomething, false);
    $t.cast(asc2, $g.generic.ASomething, false);
    $t.cast(bsc, $g.generic.BSomething, false);
    return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
  };
});
