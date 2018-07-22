$module('cast', function () {
  var $static = this;
  this.$class('d9f85bbe', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.Result = $t.property(function () {
      var $this = this;
      return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Result|3|aa28dc2d": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$interface('f5c8ee15', 'ISomeInterface', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.DoSomething = function (i) {
    return $t.cast(i, $g.cast.SomeClass, false).Result();
  };
  $static.TEST = function () {
    return $g.cast.DoSomething($g.cast.SomeClass.new());
  };
});
