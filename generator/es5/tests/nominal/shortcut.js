$module('shortcut', function () {
  var $static = this;
  this.$class('a8a36b32', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.Value = $t.property(function () {
      var $this = this;
      return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Value|3|54ff3ddf": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$type('7b5c3f1e', 'SomeNominalType', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $g.shortcut.SomeClass;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.DoSomething = function (sc) {
    return sc.Value();
  };
  $static.TEST = function () {
    var st;
    st = $t.fastbox($g.shortcut.SomeClass.new(), $g.shortcut.SomeNominalType);
    return $g.shortcut.DoSomething($t.unbox(st));
  };
});
