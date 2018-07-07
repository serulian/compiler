$module('constructable', function () {
  var $static = this;
  this.$class('5fd250f1', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $static.Get = function () {
      return $g.constructable.SomeClass.new();
    };
    $instance.SomeBool = $t.property(function () {
      var $this = this;
      return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Get|1|6caba86c<5fd250f1>": true,
        "SomeBool|3|0e92a8bc": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$interface('cd09b8fb', 'SomeInterface', false, '', function () {
    var $static = this;
    $static.Get = function () {
      return $g.constructable.SomeClass.new();
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Get|1|6caba86c<cd09b8fb>": true,
        "SomeBool|3|0e92a8bc": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomething = function (T) {
    var $f = function () {
      return T.Get();
    };
    return $f;
  };
  $static.TEST = function () {
    return $g.constructable.DoSomething($g.constructable.SomeInterface)().SomeBool();
  };
});
