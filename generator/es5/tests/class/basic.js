$module('basic', function () {
  var $static = this;
  this.$class('0021bde7', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.SomeInt = $t.fastbox(2, $g.________testlib.basictypes.Integer);
      instance.AnotherBool = $g.basic.CoolFunction();
      return instance;
    };
    $instance.AnotherFunction = function () {
      var $this = this;
      return;
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "AnotherFunction|2|6caba86c<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.CoolFunction = function () {
    return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.basic.SomeClass.new().AnotherBool;
  };
});
