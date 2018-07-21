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
      computed[("Something|2|cf412abd<" + $t.typeid(T)) + ">"] = true;
      return this.$cachedtypesig = computed;
    };
  });

  this.$class('92ebaa06', 'A', false, '', function () {
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

  this.$class('eff03fa4', 'B', false, '', function () {
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

  this.$interface('ee806320', 'ASomething', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Something|2|cf412abd<92ebaa06>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$interface('d938503c', 'BSomething', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Something|2|cf412abd<eff03fa4>": true,
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
    return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
  };
});
